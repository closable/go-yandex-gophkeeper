// Package main main package starts services
package main

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/closable/go-yandex-gophkeeper/cmd/gophkeeper/version"
	"github.com/closable/go-yandex-gophkeeper/internal/cliapp"
	"github.com/closable/go-yandex-gophkeeper/internal/config"
	"github.com/closable/go-yandex-gophkeeper/internal/logger"
	"github.com/closable/go-yandex-gophkeeper/internal/services"
)

func main() {

	ver := version.Get()
	fmt.Println(strings.Repeat("-", 50))
	fmt.Printf("BuildVersion:%s\nBuildDate:%s\nGitCommit:%s\nCompiler:%s\nPlatform:%s\n", ver.BuildVersion, ver.BuildDate, ver.GitCommit, ver.Compiler, ver.Platform)
	fmt.Println(strings.Repeat("-", 50))

	cfg := config.LoadConfig()
	logger := logger.NewLogger()

	interrupt := make(chan os.Signal, 1)
	shutdownSignals := []os.Signal{
		os.Interrupt,
		syscall.SIGTERM,
		syscall.SIGINT,
		syscall.SIGQUIT,
	}
	signal.Notify(interrupt, shutdownSignals...)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	defer logger.Sync()
	// 	если в конфиге переданы логие/пароль, то стартует режим CLI
	err := cliapp.CliAppRun(cfg.Login, cfg.Password, cfg.DSN)
	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Println("Wrong login or password. CLI was skipped!")
		}
		fmt.Println("Starting servers ...")
	}

	srv, err := services.New(cfg.DSN, cfg.ServerAddress, cfg.FileServerAddress, logger)
	if err != nil {
		panic(err)
	}

	go func() {
		<-interrupt
		srv.Shutdown(ctx)
		close(interrupt)
	}()

	srv.Run()

	<-interrupt
	fmt.Println("\nServers are shutdown gracefully !")

}
