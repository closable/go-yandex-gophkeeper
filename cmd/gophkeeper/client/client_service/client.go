package client

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"runtime"
	"strconv"
	"strings"

	"github.com/closable/go-yandex-gophkeeper/internal/cliapp"
	pb "github.com/closable/go-yandex-gophkeeper/internal/services/proto"
	"github.com/closable/go-yandex-gophkeeper/internal/store"
	"github.com/closable/go-yandex-gophkeeper/internal/utils"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

//var BuildVersion, BuildTime, buildCommit = "N/A", "N/A", "N/A"

type GKClient struct {
	Token      string
	BatchSize  int
	Client     pb.GophKeeperClient
	FileClient pb.FilseServiceClient
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
			v := cliapp.DigInput(3, marks, caret)
			if len(v) != 3 {
				return errors.New("invalid data input")
			}

			dataType, err := strconv.Atoi(v[0])
			if err != nil {
				continue
			}
			switch dataType {
			// only one file
			case 3:
				ctx, cancel := context.WithCancel(context.Background())
				defer cancel()
				err := client.UploadFile(ctx, cancel, dataType, v[2], v[1], false, 0)
				if err != nil {
					fmt.Println("При загрузке возникла ошибка ", err)
					continue
				}
			// folder
			case 4:
				ctx, cancel := context.WithCancel(context.Background())
				defer cancel()
				err := client.UploadFile(ctx, cancel, dataType, v[2], v[1], true, 0)
				if err != nil {
					fmt.Println("При загрузке возникла ошибка ", err)
					continue
				}
			default:
				err = client.AddItem(dataType, v[2], v[1])
				if err != nil {
					continue
				}
			}
		case fmt.Sprintf("d%s", caret):
			marks := []string{"ИД : "}
			v := cliapp.DigInput(1, marks, caret)
			if len(v) != 1 {
				return errors.New("invalid data input")
			}
			id, err := strconv.Atoi(v[0])
			if err != nil {
				continue
			}
			err = client.DeleteItem(id)
			if err != nil {
				fmt.Println(err)
				return err
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
			v := cliapp.DigInput(3, marks, caret)
			if len(v) != 3 {
				return errors.New("invalid data input")
			}
			dataType, err := strconv.Atoi(v[0])
			if err != nil {
				continue
			}

			id, err := strconv.Atoi(v[1])
			if err != nil {
				continue
			}
			switch dataType {
			case 3:
				ctx, cancel := context.WithCancel(context.Background())
				defer cancel()
				err := client.UploadFile(ctx, cancel, dataType, v[2], "", false, id)
				if err != nil {
					fmt.Println("При загрузке возникла ошибка ", err)
					continue
				}
			case 4:
				ctx, cancel := context.WithCancel(context.Background())
				defer cancel()
				err := client.UploadFile(ctx, cancel, dataType, v[2], "", true, id)
				if err != nil {
					fmt.Println("При загрузке возникла ошибка ", err)
					continue
				}
			default:
				err = client.UpdateItem(id, v[2])
				if err != nil {
					fmt.Println(err)
					return err
				}
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
		return err
	}
	c.Token = response.Token
	return nil
}

// ListItems list items
func (c *GKClient) ListItems(decripted bool) error {
	req := &pb.ListItemsRequest{
		Decrypted: decripted,
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

func (c *GKClient) UploadFile(ctx context.Context, cancel context.CancelFunc, dataType int, path, marker string, folder bool, dataID int) error {
	destantion := path
	if folder {
		des, err := utils.ZipFolder(path)
		if err != nil {
			fmt.Println(err)
			return err
		}
		destantion = des
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
			Chunk: chunk, Token: c.Token, DataType: int32(dataType), Name: mark, Data: path, DataID: int32(dataID),
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
