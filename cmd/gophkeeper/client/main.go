// Independent clent service
package main

import (
	"fmt"
	"log"
	"strings"
	"time"

	client "github.com/closable/go-yandex-gophkeeper/cmd/gophkeeper/client/client_service"
	"github.com/closable/go-yandex-gophkeeper/cmd/gophkeeper/version"
	"github.com/closable/go-yandex-gophkeeper/internal/config"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	cfg := config.LoadConfig()
	conn, err := grpc.NewClient(cfg.ServerAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}

	fileConn, err := grpc.NewClient(cfg.FileServerAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}

	defer conn.Close()

	c := client.New(conn, fileConn, cfg.MinioServerAddress)
	ticker := time.NewTicker(30 * time.Second)
	done := make(chan bool)
	go func() {
		client := c
		for {
			select {
			case <-done:
				return
			case <-ticker.C:
				// start sync server status
				err := client.Health()
				if err != nil {
					c.Offline = true
				} else {
					c.Offline = false
				}

			}
		}
	}()

	ver := version.Get()
	fmt.Println(strings.Repeat("-", 50))
	fmt.Printf("BuildVersion:%s\nBuildDate:%s\nGitCommit:%s\nCompiler:%s\nPlatform:%s\n", ver.BuildVersion, ver.BuildDate, ver.GitCommit, ver.Compiler, ver.Platform)
	fmt.Println(strings.Repeat("-", 50))
	if cfg.CLI {
		c.Run()
	} else {
		c.TUI()
	}
}
