// Independent clent service
package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	client "github.com/closable/go-yandex-gophkeeper/cmd/gophkeeper/client/client_service"
	"github.com/closable/go-yandex-gophkeeper/cmd/gophkeeper/version"
	"github.com/closable/go-yandex-gophkeeper/internal/config"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	cfg := config.LoadConfig()
	// s3 config
	err := godotenv.Load()
	if err != nil {
		/*
			example .env file in root folder app
				S3_BUCKET=gophkeeper
				S3_ACCESS_KEY=DNwRXfu3SAqGxZBtqLTi
				S3_SECRET_KEY=saC76TvvR5P7calPgkhvMfxO3HE68OtfaFYt1HYb
		*/
		log.Fatal("Error loading S3 environment .env file")
	}
	s3Bucket := os.Getenv("S3_BUCKET")
	s3AccessKey := os.Getenv("S3_ACCESS_KEY")
	s3SecretKey := os.Getenv("S3_SECRET_KEY")

	conn, err := grpc.NewClient(cfg.ServerAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}

	fileConn, err := grpc.NewClient(cfg.FileServerAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}

	defer conn.Close()

	c := client.New(conn, fileConn, cfg.MinioServerAddress, s3Bucket, s3AccessKey, s3SecretKey)
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
