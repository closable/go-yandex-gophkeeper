// Package cliapp create CLI and start CLI server mode
package cliapp

import (
	"bufio"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"strconv"

	"github.com/closable/go-yandex-gophkeeper/internal/store"
	"github.com/closable/go-yandex-gophkeeper/internal/utils"
)

// CliApp main structure server CLI
type CliApp struct {
	user  *store.UserDetail
	store CliStorager
}

// CliStorager интерфейс
type CliStorager interface {
	AddItem(userId, dataType int, data, name string) error
	GetUserInfo(login, password string) (*store.UserDetail, error)
	ListItems(userId int) ([]store.RowItem, error)
	UpdateItem(userId, dataId int, data string) error
	DeleteItem(userId, dataId int) error
}

// CliAppRun run cli app
func CliAppRun(login, pass, DSN string) error {
	var st CliStorager

	st, err := store.New(DSN)
	if err != nil {
		panic(err)
	}

	usr, err := st.GetUserInfo(login, pass)
	if err != nil {
		if err != sql.ErrNoRows {
			fmt.Println(err)
		}
		return err
	}

	uc := &CliApp{
		user:  usr,
		store: st,
	}
	reader := bufio.NewReader(os.Stdin)
	CliHelp()
	for {
		fmt.Print("Enter command: > ")
		text, _ := reader.ReadString('\n')

		switch c := text; c {
		case "a\n":
			marks := []string{"Тип 1-Текст 2-Ключ/Значение 3-Файл 4-Папка : ", "Метка (для файла оставлять пустой) : ", "Данные (для файла полный путь) : "}
			v := DigInput(3, marks)
			if len(v) != 3 {
				return errors.New("invalid data input")
			}

			data_type, err := strconv.Atoi(v[0])
			if err != nil {
				continue
			}

			err = uc.addItem(data_type, v[2], v[1])
			if err != nil {
				continue
			}
		case "z\n":
			marks := []string{"Метка (для файла оставлять пустой) : ", "Папка/Файл (для полный путь) : "}
			v := DigInput(2, marks)
			err := uc.zipItem(v[1])
			if err != nil {
				fmt.Println(err)
				continue
			}
		case "d\n":
			marks := []string{"ИД : "}
			v := DigInput(1, marks)
			if len(v) != 1 {
				return errors.New("invalid data input")
			}
			id, err := strconv.Atoi(v[0])
			if err != nil {
				continue
			}
			err = uc.deleteItem(id)
			if err != nil {
				return err
			}
		case "u\n":
			marks := []string{"ИД : ", "Данные :"}
			v := DigInput(1, marks)
			if len(v) != 2 {
				return errors.New("invalid data input")
			}
			id, err := strconv.Atoi(v[0])
			if err != nil {
				continue
			}
			err = uc.updateItem(id, v[1])
			if err != nil {
				return err
			}
		case "h\n":
			CliHelp()
		case "l\n":
			uc.listItems(false)
		case "p\n":
			uc.listItems(true)
		case "q\n":
			fmt.Println("Работа завершена!")
			return nil
		}
	}
}

// cliHelp cli help
func CliHelp() {
	var commands = make(map[string]string)
	commands["a"] = "Добавление"
	commands["d"] = "Удаление"
	commands["u"] = "Обновление"
	commands["z"] = "Сжать файл/папка"

	commands["h"] = "Помощь"
	commands["l"] = "Просмотреть"
	commands["p"] = "Показать пароли"
	commands["q"] = "Выход"

	commands["r"] = "Регистрация"
	commands["k"] = "Аутентификация(клиент)"

	for k, v := range commands {
		fmt.Printf("%s %s \n", k, v)
	}
}

// listItems cli list user data
func (uc *CliApp) listItems(enc bool) error {
	data, err := uc.store.ListItems(uc.user.UserID)
	if err != nil {
		return err
	}

	utils.OutputListCli(data, enc, uc.user.KeyString)
	return nil
}

// addItem cli add data
func (uc *CliApp) addItem(dataType int, data, name string) error {
	var encData string
	mark := name
	switch dataType {
	case 3:
		file, err := utils.GetFileData(data)
		if err != nil {
			fmt.Println(err)
			return err
		}
		mark = data
		encData = utils.Encrypt(uc.user.KeyString, string(file))

		// store enc data into bin file
		pathBin := utils.MakePathFile(data, "bin")
		err = utils.StoreFileData(pathBin, encData)
		if err != nil {
			fmt.Println(err)
			return nil
		}
	case 4:
		// create Zip arhive
		zipPath, err := utils.ZipFolder(data)
		if err != nil {
			fmt.Println(err)
			return err
		}
		// read data from arhive
		file, err := utils.GetFileData(zipPath)
		if err != nil {
			fmt.Println(err)
			return err
		}
		mark = data
		// encrypt data
		encData = utils.Encrypt(uc.user.KeyString, string(file))

		// store enc data into bin file
		pathBin := utils.MakePathFile(data, "bin")
		err = utils.StoreFileData(pathBin, encData)
		if err != nil {
			fmt.Println(err)
			return nil
		}
		// remove zip arhive
		// err := os.RemoveAll("directoryname") is needed
		err = os.Remove(zipPath)
		if err != nil {
			fmt.Println(err)
			return nil
		}
		fmt.Println("Операция завершена")
	default:
		encData = utils.Encrypt(uc.user.KeyString, data)
	}

	err := uc.store.AddItem(uc.user.UserID, dataType, encData, mark)
	if err != nil {
		fmt.Println("Ошибка при добавлении данных", err)
		return err
	}
	fmt.Println("Данные успешно добавлены!")
	return nil
}

// updateItem update item info
func (uc *CliApp) updateItem(id int, data string) error {
	encData := utils.Encrypt(uc.user.KeyString, data)
	err := uc.store.UpdateItem(uc.user.UserID, id, encData)
	if err != nil {
		fmt.Println("Ошибка при обновлении данных!", err)
		return err
	}
	fmt.Println("Данные успешно обновлены!")
	return nil
}

// deleteItem delete selected row
func (uc *CliApp) deleteItem(id int) error {
	err := uc.store.DeleteItem(uc.user.UserID, id)
	if err != nil {
		return err
	}

	fmt.Println("Данные удалены!")
	return nil
}

// zipItem create zip arhive a folder
func (uc *CliApp) zipItem(path string) error {
	zipPath, err := utils.ZipFolder(path)
	if err != nil {
		fmt.Println(err)
		return err
	}

	err = uc.addItem(4, zipPath, zipPath)
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

// digInput helper function for detail input
func DigInput(n int, m []string) []string {

	reader := bufio.NewReader(os.Stdin)
	result := make([]string, 0)

	for {
		fmt.Print(m[len(result)])
		text, _ := reader.ReadString('\n')
		v := text[:len(text)-1]
		result = append(result, v)

		if len(result) == n {
			return result
		}

		if text == "q\n" {
			break
		}
	}
	return result
}
