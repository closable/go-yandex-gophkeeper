package cliapp

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/closable/go-yandex-gophkeeper/internal/store"
	"github.com/closable/go-yandex-gophkeeper/internal/utils"
)

type CliApp struct {
	user *store.UserDetail
	db   *store.Store
}

// CliAppRun run cli app
func CliAppRun(u *store.UserDetail, db *store.Store) error {
	uc := &CliApp{
		user: u,
		db:   db,
	}
	reader := bufio.NewReader(os.Stdin)
	cliHelp()
	for {
		fmt.Print("Enter command: > ")
		text, _ := reader.ReadString('\n')

		switch c := text; c {
		case "a\n":
			marks := []string{"Тип 1-Текст 2-Ключ/Значение 3-Файл 4-Папка : ", "Метка (для файла оставлять пустой) : ", "Данные (для файла полный путь) : "}
			v := digInput(3, marks)
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

		case "d\n":
			marks := []string{"ИД : "}
			v := digInput(1, marks)
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
			v := digInput(1, marks)
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
			cliHelp()
		case "l\n":
			uc.listItems(false)
		case "p\n":
			uc.listItems(true)
		case "q\n":
			fmt.Println("Работа завершена")
			return nil
		}
	}
}

// cliHelp cli help
func cliHelp() {
	var commands = make(map[string]string)
	commands["a"] = "Добавление"
	commands["d"] = "Удалиние"
	commands["u"] = "Обновление"

	commands["h"] = "Помощь"
	commands["l"] = "Просмотреть"
	commands["p"] = "Показать пароли"
	commands["q"] = "Выход"

	for k, v := range commands {
		fmt.Printf("%s %s \n", k, v)
	}
}

// listItems cli list user data
func (uc *CliApp) listItems(enc bool) error {
	data, err := uc.db.ListItems(uc.user.UserID)
	if err != nil {
		return err
	}
	fmt.Println(strings.Repeat("_", 123))
	fmt.Printf("| %3s | %-15s | %-25s | %-55s| %-10s |\n", "ИД", "Тип", "Метка", "Данные", "Размер,байт")
	fmt.Println(strings.Repeat("-", 123))

	for _, v := range data {
		data := v.EncData
		name := catStringData(v.Name, 20)
		if enc && v.DataType < 3 {
			data = utils.Decrypt(uc.user.KeyString, v.EncData)
		}
		data = catStringData(data, 50)
		fmt.Printf("| %3d | %-15s | %-25s | %-55s| %11d |\n", v.Id, v.Type, name, data, v.Length)
	}
	fmt.Println(strings.Repeat("-", 123))
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
		pathBin := utils.MakePathBinFile(data)
		err = utils.StoreFileData(pathBin, encData)
		if err != nil {
			fmt.Println(err)
			return nil
		}

	default:
		encData = utils.Encrypt(uc.user.KeyString, data)
	}

	err := uc.db.AddItem(uc.user.UserID, dataType, encData, mark)
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
	err := uc.db.UpdateItem(uc.user.UserID, id, encData)
	if err != nil {
		fmt.Println("Ошибка при обновлении данных!", err)
		return err
	}
	fmt.Println("Данные успешно обновлены!")
	return nil
}

// deleteItem delete selected row
func (uc *CliApp) deleteItem(id int) error {
	err := uc.db.DeleteItem(uc.user.UserID, id)
	if err != nil {
		return err
	}

	fmt.Println("Данные удалены!")
	return nil
}

// digInput helper function for detail input
func digInput(n int, m []string) []string {

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

// catStringData helper for formatting string
func catStringData(s string, n int) string {
	if len([]byte(s)) > n {
		return s[:n+1] + "..."
	}
	return s
}
