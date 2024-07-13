package main

import (
	"fmt"

	"github.com/closable/go-yandex-gophkeeper/internal/cliapp"
	"github.com/closable/go-yandex-gophkeeper/internal/config"
	"github.com/closable/go-yandex-gophkeeper/internal/store"
)

var buildVersion, buildDate, buildCommit = "N/A", "N/A", "N/A"

func main() {
	// go run -ldflags "-X main.buildVersion=v1.0.1 -X 'main.buildDate=$(date +'%Y/%m/%d %H:%M:%S')'" cmd/gophkeeper/main.go
	// go build -ldflags "-X main.buildVersion=v1.0.1 -X 'main.buildDate=$(date +'%Y/%m/%d %H:%M:%S')'" cmd/gophkeeper/main.go
	// start bin file -> ./main
	fmt.Printf("Build version:%s\nBuild date:%s\nBuild commit:%s\n", buildVersion, buildDate, buildCommit)

	// originalText := "Hello world"
	// keyStr, err := utils.CryptoSeq(secretPhrase)
	// if err != nil {
	// 	panic(err)
	// }
	// cryptoText := utils.Encrypt(keyStr, originalText)
	// fmt.Println(cryptoText)

	// text := utils.Decrypt(keyStr, cryptoText)
	// fmt.Println(text)

	// // ----
	//data, _ := os.ReadFile("book.pdf")
	// cryptoText = utils.Encrypt(keyStr, string(data))
	// //fmt.Println(cryptoText)

	// file, _ := os.Create("/tmp/book.bin")
	// defer file.Close()
	// _, err = file.WriteString(cryptoText)
	// if err != nil {
	// 	fmt.Println(err)
	// }

	// err = os.Remove("book.pdf")
	// if err != nil {
	// 	fmt.Println(err)
	// }

	// text = utils.Decrypt(keyStr, cryptoText)
	// file, _ = os.Create("book_decripted.pdf")
	// defer file.Close()
	// _, err = file.WriteString(text)
	// if err != nil {
	// 	fmt.Println(err)
	// }

	// fmt.Println(text)

	cfg := config.LoadConfig()

	db, err := store.New(cfg.DSN)
	if err != nil {
		panic(err)
	}

	fmt.Println(cfg.Login, cfg.Password)
	usr, err := db.GetUserInfo(cfg.Login)
	if err != nil {
		fmt.Println(err)
	}

	_ = cliapp.CliAppRun(usr, db)

	// phrase := "test"
	// keyStr, _ := utils.CryptoSeq(phrase)

	// user := utils.Encrypt(keyStr, "kapa01")
	// pass := utils.Encrypt(keyStr, "130301")

	// usr, err := db.CreateUser(user, pass, keyStr)
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// fmt.Println("!!!", usr)

	// usr, err := db.GetUserInfo("kapa")
	// fmt.Println(usr, err)

	// data := utils.Encrypt(usr.KeyString, `{"login": "kapa", "password": "1303"}`)
	// err = db.AddItem(usr.UserID, 2, data, "secret")
	// fmt.Println(err)

	//list
	// items, err := db.ListItems(usr.UserID)
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// for _, v := range items {
	// 	fmt.Println(v.Id, v.Type, v.Name, utils.Decrypt(usr.KeyString, v.EncData))
	// }

	// delete
	// db.DeleteItem(usr.UserID, 6)
	// items, err = db.ListItems(usr.UserID)
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// for _, v := range items {
	// 	fmt.Println(v.Id, v.Type, v.Name, utils.Decrypt(usr.KeyString, v.EncData))
	// }

	//update
	//data := utils.Encrypt(usr.KeyString, `{"login": "kapa", "password": "kapran1303"}`)
	//db.UpdateItem(usr.UserID, 9, data)
}
