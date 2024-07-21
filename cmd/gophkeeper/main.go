package main

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/closable/go-yandex-gophkeeper/cmd/gophkeeper/version"
	"github.com/closable/go-yandex-gophkeeper/internal/cliapp"
	"github.com/closable/go-yandex-gophkeeper/internal/config"
	"github.com/closable/go-yandex-gophkeeper/internal/services"
)

func main() {

	ver := version.Get()
	fmt.Println(strings.Repeat("-", 50))
	fmt.Printf("BuildVersion:%s\nBuildDate:%s\nGitCommit:%s\nCompiler:%s\nPlatform:%s\n", ver.BuildVersion, ver.BuildDate, ver.GitCommit, ver.Compiler, ver.Platform)
	fmt.Println(strings.Repeat("-", 50))

	cfg := config.LoadConfig()

	// 	если в конфиге переданы логие/пароль, то стартует режим CLI
	err := cliapp.CliAppRun(cfg.Login, cfg.Password, cfg.DSN)
	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Println("Wrong login or password. CLI was skipped!")
		}
		fmt.Println("Starting server mode ...")
	}

	srv, err := services.New(cfg.DSN, cfg.ServerAddress, cfg.FileServerAddress)
	if err != nil {
		panic(err)
	}
	srv.Run()

}
