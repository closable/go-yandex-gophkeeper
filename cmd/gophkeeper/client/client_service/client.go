// Package client  realize clinet for work with server serice
package client

import (
	"bufio"
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"runtime"
	"strconv"
	"strings"
	"sync"

	"github.com/charmbracelet/bubbles/table"
	"github.com/closable/go-yandex-gophkeeper/cmd/gophkeeper/client/tui"
	"github.com/closable/go-yandex-gophkeeper/cmd/gophkeeper/client/tui/models"
	tuitable "github.com/closable/go-yandex-gophkeeper/cmd/gophkeeper/client/tui/table"
	"github.com/closable/go-yandex-gophkeeper/cmd/gophkeeper/client/tui/textinput"
	"github.com/closable/go-yandex-gophkeeper/internal/cliapp"
	errs "github.com/closable/go-yandex-gophkeeper/internal/errors"
	pb "github.com/closable/go-yandex-gophkeeper/internal/services/proto"
	"github.com/closable/go-yandex-gophkeeper/internal/store"
	"github.com/closable/go-yandex-gophkeeper/internal/utils"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type (
	// GKClient client structure
	GKClient struct {
		Token      string
		KeyString  string
		BatchSize  int
		Client     pb.GophKeeperClient
		FileClient pb.FilseServiceClient
		Cache      LocalCache
		Offline    bool
	}
	// LocalCache cache structure
	LocalCache struct {
		Store map[int]store.RowItem
		mu    sync.Mutex
	}
)

var (
	columns = []table.Column{
		{Title: "ИД", Width: 10},
		{Title: "Тип", Width: 20},
		{Title: "Метка", Width: 40},
		{Title: "Данные", Width: 40},
		{Title: "Размер,байт", Width: 10},
	}
)

// NewLocalCache new cache instance
func NewLocalCache() *LocalCache {
	var st = make(map[int]store.RowItem)

	return &LocalCache{
		Store: st,
	}
}

// CacheDecode decode cache file if exists
func (lc *LocalCache) CacheDecode(key string) error {
	f, err := os.Open("./cache")
	if err != nil {
		fmt.Println("Файл не найден! ", err)
		return err
	}
	defer f.Close()

	stat, err := f.Stat()
	if err != nil {
		fmt.Println("Файл кэша не найден! ", err)
		return err
	}

	b := make([]byte, stat.Size())
	_, err = f.Read(b)
	if err != nil {
		fmt.Println("Ошибка чтения файла ", err)
		return err
	}

	data, err := base64.StdEncoding.DecodeString(string(b))
	if err != nil {
		fmt.Println("Ошибка декодирования ", err)
		return err
	}

	encData := utils.Encrypt(key, string(data))

	var st = make(map[int]store.RowItem)

	err = json.Unmarshal([]byte(encData), &st)
	if err != nil {
		return err
	}

	lc.Store = st

	return nil
}

// Sync sync data cache
func (lc *LocalCache) Sync(r []store.RowItem) error {
	// update all records
	for _, v := range r {
		lc.Add(v)
	}
	return nil
}

// ToFile store cache data to local file
func (lc *LocalCache) ToFile(key string) error {
	s, err := json.Marshal(lc.Store)
	if err != nil {
		fmt.Println(err)
		return err
	}
	encData := utils.Encrypt(key, string(s))
	dec := base64.StdEncoding.EncodeToString([]byte(encData))
	utils.StoreFileData("./cache", dec)
	fmt.Println("Данные сохранены!")

	return nil
}

// Clear func clear data in cache
func (lc *LocalCache) Clear() {
	lc.Store = make(map[int]store.RowItem)
}

// Add func add data to cache
func (lc *LocalCache) Add(r store.RowItem) {
	lc.mu.Lock()
	lc.Store[r.Id] = store.RowItem{
		Id:        r.Id,
		Type:      r.Type,
		Name:      r.Name,
		IsRestore: r.IsRestore,
		EncData:   r.EncData,
		Length:    r.Length,
		DataType:  r.DataType,
	}
	lc.mu.Unlock()
}

// List cache data
func (lc *LocalCache) List() []store.RowItem {
	res := make([]store.RowItem, 0)
	for _, v := range lc.Store {
		res = append(res, v)
	}
	return res
}

// Run start CLI mode
func (c *GKClient) Run() error {
	cliapp.CliHelp()

	if len(c.Token) < 10 {
		fmt.Println(strings.Repeat("-", 50))
		fmt.Println("Аутентификаця не выполнена, начните с нее!")
		fmt.Println(strings.Repeat("-", 50))
	}
	reader := bufio.NewReader(os.Stdin)
	client := c

	caret := "\n"
	if runtime.GOOS == "windows" {
		caret = "\r\n"
	}

	for {
		fmt.Print("Enter command: > ")
		text, _ := reader.ReadString('\n')

		switch c := text; c {
		case fmt.Sprintf("h%s", caret):
			cliapp.CliHelp()
		case fmt.Sprintf("q%s", caret):
			fmt.Println("Работа завершена!")
			return nil
		case fmt.Sprintf("g%s", caret):
			marks := []string{"ИД : "}
			v := cliapp.DigInput(1, marks, caret)
			if len(v) != 1 {
				return errors.New("invalid data input")
			}
			id, err := strconv.Atoi(v[0])
			if err != nil {
				continue
			}
			err = client.DownloadFile(id)
			if err != nil {
				fmt.Println(err)
			}

		case fmt.Sprintf("a%s", caret):
			marks := []string{"Тип 1-Текст 2-Ключ/Значение 3-Файл 4-Папка : ", "Метка (для файла оставлять пустой) : ", "Данные (для файла полный путь) : "}
			data := cliapp.DigInput(3, marks, caret)
			if len(data) != 3 {
				return errors.New("invalid data input")
			}
			err := client.add(data)
			if err != nil {
				fmt.Println(err)
				continue
			}
		case fmt.Sprintf("d%s", caret):
			marks := []string{"ИД : "}
			data := cliapp.DigInput(1, marks, caret)
			if len(data) != 1 {
				return errors.New("invalid data input")
			}
			err := client.delete(data)
			if err != nil {
				fmt.Println(err)
				continue
			}
		case fmt.Sprintf("r%s", caret):
			marks := []string{"Логин : ", "Пароль :", "Сектретная фраза (можно оставить пустым) :"}
			v := cliapp.DigInput(3, marks, caret)
			if len(v) != 3 {
				return errors.New("invalid data input")
			}

			err := client.CreateUser(v[0], v[1], v[2])
			if err != nil {
				fmt.Println(err)
				return err
			}
		case fmt.Sprintf("u%s", caret):
			marks := []string{"Тип 1-Текст 2-Ключ/Значение 3-Файл 4-Папка : ", "ИД : ", "Данные :"}
			data := cliapp.DigInput(3, marks, caret)
			if len(data) != 3 {
				return errors.New("invalid data input")
			}

			err := client.update(data)
			if err != nil {
				fmt.Println(err)
				continue
			}
		case fmt.Sprintf("l%s", caret):
			client.ListItems(false)
		case fmt.Sprintf("p%s", caret):
			client.ListItems(true)
		case fmt.Sprintf("k%s", caret):
			marks := []string{"Логин : ", "Пароль : "}
			v := cliapp.DigInput(2, marks, caret)
			if len(v) != 2 {
				return errors.New("invalid data input")
			}
			err := client.Login(v[0], v[1])
			if err != nil {
				fmt.Println(err)
				continue
			}
		}
	}
}

// TUI start TUI mode
func (c *GKClient) TUI() error {
	client := c

	for {
		choice, err := tui.MainMenu(c.Token)
		if err != nil || choice < 0 {
			break
		}

		switch choice {
		case 0: // create
			if client.Offline {
				continue
			}
			m := make([]models.TuiModelText, 0)
			m = append(m, models.TuiModelText{Label: "Login", IsEcho: false, CharLimit: 64})
			m = append(m, models.TuiModelText{Label: "Password", IsEcho: true, CharLimit: 32})
			m = append(m, models.TuiModelText{Label: "Secret phrase", IsEcho: false, CharLimit: 32})
			data, err := textinput.TUItext(m)
			if err != nil {
				continue
			}
			err = client.CreateUser(data[0], data[1], data[2])
			if err != nil {
				fmt.Println(err)
				client.Offline = true
				continue
			}
		case 1: // login
			if client.Offline {
				continue
			}
			m := make([]models.TuiModelText, 0)
			m = append(m, models.TuiModelText{Label: "Login", IsEcho: false, CharLimit: 64})
			m = append(m, models.TuiModelText{Label: "Password", IsEcho: true, CharLimit: 32})
			data, err := textinput.TUItext(m)
			if err != nil {
				fmt.Printf("%v %v", errs.ErrorLogin.Error(), err)
				continue
			}
			err = client.Login(data[0], data[1])
			if err != nil {
				fmt.Printf("%v %v", errs.ErrorLogin.Error(), err)
				client.Offline = true
				continue
			}

			client.fullSyncCache(false)
		case 2: // List
			var data []store.RowItem
			data, err = client.ListItemsData(false)
			if err != nil {
				fmt.Println(err)
				fmt.Println("Warning ! off-line mode")
				client.Offline = true
				data = client.Cache.List()
			}
			tuitable.TUItable(columns, dataToRowsTUI(data))
		case 3: // List decrypted
			data, err := client.ListItemsData(true)
			if err != nil {
				fmt.Println("Warning ! off-line mode")
				data = client.Cache.List()
			}
			tuitable.TUItable(columns, dataToRowsTUI(data))
		case 4: // Add
			if len(c.Token) < 10 || client.Offline {
				continue
			}
			m := make([]models.TuiModelText, 0)
			m = append(m, models.TuiModelText{Label: "Тип 1-Текст 2-Ключ/Значение 3-Файл 4-Папка", IsEcho: false, CharLimit: 32})
			m = append(m, models.TuiModelText{Label: "Метка (для файла оставлять пустой)", IsEcho: false, CharLimit: 32})
			m = append(m, models.TuiModelText{Label: "Данные (для файла полный путь)", IsEcho: false, CharLimit: 255})
			data, err := textinput.TUItext(m)
			if err != nil {
				continue
			}

			err = client.add(data)
			if err != nil {
				fmt.Println(err)
				continue
			}
			client.fullSyncCache(true)
		case 5: // Update
			if len(c.Token) < 10 || client.Offline {
				continue
			}
			m := make([]models.TuiModelText, 0)
			m = append(m, models.TuiModelText{Label: "Тип 1-Текст 2-Ключ/Значение 3-Файл 4-Папка", IsEcho: false, CharLimit: 32})
			m = append(m, models.TuiModelText{Label: "ИД", IsEcho: false, CharLimit: 32})
			m = append(m, models.TuiModelText{Label: "Данные (для файла полный путь)", IsEcho: false, CharLimit: 32})
			data, err := textinput.TUItext(m)
			if err != nil {
				continue
			}

			err = c.update(data)
			if err != nil {
				fmt.Println(err)
				continue
			}
			client.fullSyncCache(true)
		case 6: //Get file/folder
			if len(c.Token) < 10 || client.Offline {
				continue
			}
			m := make([]models.TuiModelText, 0)
			m = append(m, models.TuiModelText{Label: "ИД", IsEcho: false, CharLimit: 64})

			data, err := textinput.TUItext(m)
			if err != nil {
				continue
			}

			id, err := strconv.Atoi(data[0])
			if err != nil {
				continue
			}
			err = client.DownloadFile(id)
			if err != nil {
				fmt.Println(err)
			}
		case 7: // Delete
			if len(c.Token) < 10 || client.Offline {
				continue
			}
			m := make([]models.TuiModelText, 0)
			m = append(m, models.TuiModelText{Label: "ИД", IsEcho: false, CharLimit: 64})

			data, err := textinput.TUItext(m)
			if err != nil {
				continue
			}

			err = client.delete(data)
			if err != nil {
				fmt.Println(err)
				continue
			}
			client.fullSyncCache(true)
		case 9:
			err := client.Cache.ToFile(c.KeyString)
			if err != nil {
				fmt.Println(err)
			}
			fmt.Println("Работа завершена!")
			return nil
		}
	}
	return nil
}

// dataToRowsTUI func transformation types of data
func dataToRowsTUI(d []store.RowItem) []table.Row {
	rows := make([]table.Row, 0)
	for _, v := range d {
		rows = append(rows, table.Row{
			fmt.Sprintf("%d", v.Id), v.Type, v.Name, v.EncData, fmt.Sprintf("%d", v.Length),
		})
	}
	return rows
}

// fullSyncCache func for sync cache and server data
func (c *GKClient) fullSyncCache(clearBefore bool) {
	client := c
	if clearBefore {
		c.Cache.Clear()
	}

	// get all records from db
	rows, err := client.ListItemsData(true)
	if err != nil {
		client.Offline = true
		fmt.Println(err)
	}
	// sync with local cache
	err = c.Cache.Sync(rows)
	if err != nil {
		fmt.Println(err)
	}

}

// add helper function
func (c *GKClient) add(d []string) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	client := c
	dataType, err := strconv.Atoi(d[0])
	if err != nil {
		return err
	}
	switch dataType {
	case 2:
		type keyValue struct {
			Key   string
			Value string
		}

		m := make([]models.TuiModelText, 0)
		m = append(m, models.TuiModelText{Label: "Наименование ключа", IsEcho: false, CharLimit: 255})
		m = append(m, models.TuiModelText{Label: "Значение ключа", IsEcho: false, CharLimit: 255})
		data, err := textinput.TUItext(m)
		if err != nil {
			return err
		}

		v := &keyValue{Key: data[0], Value: data[1]}

		js, err := json.Marshal(v)
		if err != nil {
			return err
		}

		err = client.AddItem(dataType, string(js), d[1])
		if err != nil {
			return err
		}

	// only one file
	case 3:
		err := client.UploadFile(ctx, cancel, dataType, d[2], d[1], false, 0)
		if err != nil {
			return err
		}
	// folder
	case 4:
		err := client.UploadFile(ctx, cancel, dataType, d[2], d[1], true, 0)
		if err != nil {
			return err
		}
	default:
		err = client.AddItem(dataType, d[2], d[1])
		if err != nil {
			return err
		}
	}
	return nil
}

// update helper function
func (c *GKClient) update(d []string) error {
	client := c
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	dataType, err := strconv.Atoi(d[0])
	if err != nil {
		return err
	}

	id, err := strconv.Atoi(d[1])
	if err != nil {
		return err
	}
	switch dataType {
	case 3:
		err := client.UploadFile(ctx, cancel, dataType, d[2], "", false, id)
		if err != nil {
			return err
		}
	case 4:
		err := client.UploadFile(ctx, cancel, dataType, d[2], "", true, id)
		if err != nil {
			return err
		}
	default:
		err = client.UpdateItem(id, d[2])
		if err != nil {
			return err
		}
	}
	return nil
}

// delete helper function
func (c *GKClient) delete(d []string) error {
	client := c
	id, err := strconv.Atoi(d[0])
	if err != nil {
		return err
	}
	err = client.DeleteItem(id)
	if err != nil {
		return err
	}
	return nil
}

// Health check health status
func (c *GKClient) Health() error {

	md := metadata.New(map[string]string{"authorization": fmt.Sprintf("Bearer %s", c.Token)})
	ctx := metadata.NewOutgoingContext(context.Background(), md)
	var header metadata.MD

	req := &pb.HealthRequest{Numb: "1"}

	_, err := c.Client.Health(ctx, req, grpc.Header(&header))
	if err != nil {
		return err
	}

	return nil
}

// Login auth user
func (c *GKClient) Login(user, pass string) error {
	req := &pb.LoginRequest{
		User: user,
		Pass: pass,
	}
	ctx := context.Background()
	response, err := c.Client.Login(ctx, req)
	if err != nil {
		c.Token = "*"
		c.KeyString = ""
		return err
	}
	c.Token = response.Token
	c.KeyString = response.KeyString

	c.Cache.CacheDecode(c.KeyString)
	return nil
}

// ListItems list items
func (c *GKClient) ListItems(decrypted bool) error {
	req := &pb.ListItemsRequest{
		Decrypted: decrypted,
	}
	md := metadata.New(map[string]string{"authorization": fmt.Sprintf("Bearer %s", c.Token)})
	ctx := metadata.NewOutgoingContext(context.Background(), md)
	var header metadata.MD

	response, err := c.Client.ListItems(ctx, req, grpc.Header(&header))
	if err != nil {
		return err
	}
	data := make([]store.RowItem, 0)
	for _, v := range response.Items {
		data = append(data, store.RowItem{
			Id:        int(v.Id),
			Type:      v.Type,
			Name:      v.Name,
			IsRestore: v.Restore,
			EncData:   v.Encdata,
			Length:    int(v.Length),
			DataType:  int(v.DataType),
		})
	}

	utils.OutputListCli(data, false, "")

	return nil
}

// ListItems list items
func (c *GKClient) ListItemsData(decrypted bool) ([]store.RowItem, error) {
	req := &pb.ListItemsRequest{
		Decrypted: decrypted,
	}
	md := metadata.New(map[string]string{"authorization": fmt.Sprintf("Bearer %s", c.Token)})
	ctx := metadata.NewOutgoingContext(context.Background(), md)
	var header metadata.MD
	data := make([]store.RowItem, 0)

	response, err := c.Client.ListItems(ctx, req, grpc.Header(&header))
	if err != nil {
		return data, err
	}

	for _, v := range response.Items {
		data = append(data, store.RowItem{
			Id:        int(v.Id),
			Type:      v.Type,
			Name:      v.Name,
			IsRestore: v.Restore,
			EncData:   v.Encdata,
			Length:    int(v.Length),
			DataType:  int(v.DataType),
		})
	}
	return data, nil
}

// AddItem func add data
func (c *GKClient) AddItem(dataType int, data, name string) error {
	req := &pb.AddItemRequest{
		DataType: int32(dataType),
		Data:     data,
		Name:     name,
	}

	md := metadata.New(map[string]string{"authorization": fmt.Sprintf("Bearer %s", c.Token)})
	ctx := metadata.NewOutgoingContext(context.Background(), md)
	var header metadata.MD

	_, err := c.Client.AddItem(ctx, req, grpc.Header(&header))
	if err != nil {
		return err
	}

	return nil
}

// UpdateItem func update stored data
func (c *GKClient) UpdateItem(dataID int, data string) error {
	req := &pb.UpdateItemRequest{
		DataID: int32(dataID),
		Data:   data,
	}
	md := metadata.New(map[string]string{"authorization": fmt.Sprintf("Bearer %s", c.Token)})
	ctx := metadata.NewOutgoingContext(context.Background(), md)
	var header metadata.MD

	_, err := c.Client.UpdateItem(ctx, req, grpc.Header(&header))
	if err != nil {
		return err
	}

	return nil
}

// DeleteItem func delete data
func (c *GKClient) DeleteItem(dataId int) error {
	req := &pb.DelItemRequest{
		DataID: int32(dataId),
	}
	md := metadata.New(map[string]string{"authorization": fmt.Sprintf("Bearer %s", c.Token)})
	ctx := metadata.NewOutgoingContext(context.Background(), md)
	var header metadata.MD
	_, err := c.Client.DelItem(ctx, req, grpc.Header(&header))
	if err != nil {
		return err
	}

	return nil
}

// CreateUser func create user
func (c *GKClient) CreateUser(login, password, secret string) error {
	req := &pb.CreateUserRequest{
		User:      login,
		Pass:      password,
		Keystring: secret,
	}
	ctx := context.Background()
	response, err := c.Client.CreateUser(ctx, req)
	if err != nil {
		c.Token = "*"
		return err
	}

	c.Token = response.Token
	return nil
}

// DownloadFile func download binary data
func (c *GKClient) DownloadFile(dataID int) error {
	req := &pb.FileDownloadRequest{
		DataID: int32(dataID),
		Token:  c.Token,
	}
	stream, err := c.FileClient.Download(context.Background(), req)
	if err != nil {
		fmt.Println(err)
		return err
	}

	file, err := os.CreateTemp("", fmt.Sprintf("tmp_%s", utils.GetRandomString(10)))
	//file, err := os.Create(fmt.Sprintf("tmp_%s", utils.GetRandomString(10)))
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer file.Close()
	filePath := ""
	dataType := 0
	for {
		chank, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			break
		}

		_, err = file.Write(chank.GetChank())
		if err != nil {
			fmt.Println(err)
			return err
		}

		if dataType == 0 {
			dataType = int(chank.GetDataType())
			filePath = chank.GetFilePath()
		}
	}

	fileDest := filePath
	if _, err := os.Stat(fileDest); !os.IsNotExist(err) && dataType == 3 {
		path := strings.Split(filePath, string(os.PathSeparator))
		filename := strings.Split(path[len(path)-1], ".")
		if len(filename) == 2 {
			path[len(path)-1] = fmt.Sprintf("%s(*).%s", filename[0], filename[1])
		} else {
			path[len(path)-1] = fmt.Sprintf("%s(*)", filename[0])
		}

		fileDest = strings.Join(path, string(os.PathSeparator))
	}
	fmt.Println("Расположение файла -", fileDest)
	if dataType == 4 {
		fileDest = fmt.Sprintf("%s.zip", fileDest)
	}

	dest, err := os.Create(fileDest)
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer dest.Close()

	_, err = file.Seek(int64(0), 0)
	if err != nil {
		fmt.Println(err)
		return err
	}

	_, err = io.Copy(dest, file)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

// UploadFile func upload binary data
func (c *GKClient) UploadFile(ctx context.Context, cancel context.CancelFunc, dataType int, path, marker string, folder bool, dataID int) error {
	destantion := path
	if folder {
		des, err := utils.ZipFolder(path)
		if err != nil {
			fmt.Println(err)
			return err
		}
		destantion = des
	} else {
		f, err := os.Open(destantion)
		if err != nil {
			return err
		}

		i, err := f.Stat()
		if err != nil {
			return err
		}
		if i.Size()/10000000 > 10 {
			des, err := utils.ZipFile(path)
			if err != nil {
				fmt.Println(err)
				return err
			}
			destantion = des
		}
		f.Close()
	}

	stream, err := c.FileClient.Upload(ctx)
	if err != nil {
		return err
	}
	mark := marker
	if len(marker) == 0 {
		mark = path
	}
	file, err := os.Open(destantion)
	if err != nil {
		return err
	}
	buf := make([]byte, c.BatchSize)
	batchNumber := 1

	for {
		num, err := file.Read(buf)
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		chunk := buf[:num]
		if err := stream.Send(&pb.FileUploadRequest{
			Chunk: chunk, Token: c.Token, DataType: int32(dataType), Name: mark, Data: destantion, DataID: int32(dataID),
		}); err != nil {
			return err
		}
		batchNumber += 1
	}

	_, err = stream.CloseAndRecv()
	if err != nil {
		stream.CloseSend()
	}

	cancel()
	return nil
}

func New(conn, fileConn *grpc.ClientConn) *GKClient {
	return &GKClient{
		Client:     pb.NewGophKeeperClient(conn),
		FileClient: pb.NewFilseServiceClient(fileConn),
		Token:      "",
		BatchSize:  1024,
		Cache:      *NewLocalCache(),
		Offline:    false,
	}
}
