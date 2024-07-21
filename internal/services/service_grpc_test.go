package services

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"testing"

	"github.com/closable/go-yandex-gophkeeper/cmd/gophkeeper/client"
	"github.com/closable/go-yandex-gophkeeper/internal/config"
	pb "github.com/closable/go-yandex-gophkeeper/internal/services/proto"
	"github.com/closable/go-yandex-gophkeeper/internal/store"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

func TestCreateUser(t *testing.T) {
	cfg := config.LoadConfig()

	srv, err := New(cfg.DSN, cfg.ServerAddress, cfg.FileServerAddress)
	if err != nil {
		panic(err)
	}

	go srv.Run()

	conn, err := grpc.NewClient(srv.ServAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}

	fileConn, err := grpc.NewClient(srv.FileServAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}

	tearDown, _ := store.New(cfg.DSN)

	defer conn.Close()

	c := &client.GKClient{
		Client:     pb.NewGophKeeperClient(conn),
		FileClient: pb.NewFilseServiceClient(fileConn),
		Token:      "",
		BatchSize:  1024,
	}

	tests := []struct {
		name     string
		login    string
		password string
	}{
		{
			name:     "create user 1",
			login:    fmt.Sprintf("user_%d", rand.Intn(100)),
			password: "test1",
		},
		{
			name:     "create user 2",
			login:    fmt.Sprintf("user_%d", rand.Intn(100)),
			password: "test2",
		},
		{
			name:     "create user 3",
			login:    fmt.Sprintf("user_%d", rand.Intn(100)),
			password: "test3",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()

			// new && login
			req := &pb.CreateUserRequest{
				User:      tt.login,
				Pass:      tt.name,
				Keystring: "",
			}
			resp, err := c.Client.CreateUser(ctx, req)
			if err != nil {
				t.Errorf("Create user error %v %v", tt.name, err)
			}

			if len(resp.Token) < 10 {
				t.Errorf("Wrong authentication %v %v %v", tt.name, resp.User, req.User)
			}

			c.Token = resp.Token
			md := metadata.New(map[string]string{"authorization": fmt.Sprintf("Bearer %s", c.Token)})
			ctx = metadata.NewOutgoingContext(context.Background(), md)
			var header metadata.MD

			if resp.User.Login != req.User {
				t.Errorf("User name must be equal %v %v %v", tt.name, resp.User, req.User)
			}

			// add item
			addReq := &pb.AddItemRequest{
				DataType: 1,
				Data:     "hello world",
				Name:     "pass phrase",
			}
			_, err = c.Client.AddItem(ctx, addReq, grpc.Header(&header))
			if err != nil {
				t.Errorf("Add data error %v %v", tt.name, err)
			}

			// output info
			listReq := &pb.ListItemsRequest{Decrypted: true}
			list, err := c.Client.ListItems(ctx, listReq, grpc.Header(&header))
			if err != nil {
				t.Errorf("Error list information %v %v", tt.name, err)
			}

			if len(list.Items) == 0 {
				t.Errorf("Length must be greater than 0 %v %v", tt.name, err)
			}

			// update
			updReq := &pb.UpdateItemRequest{
				DataID: list.Items[0].Id,
				Data:   "1234567890",
			}

			_, err = c.Client.UpdateItem(ctx, updReq, grpc.Header(&header))
			if err != nil {
				t.Errorf("Update data error %v %v", tt.name, err)
			}

			// check status update
			list, err = c.Client.ListItems(ctx, listReq, grpc.Header(&header))
			if err != nil {
				t.Errorf("Error list information after update %v %v", tt.name, err)
			}

			if list.Items[0].Encdata != "1234567890" {
				t.Errorf("Wrong information after update %v %v", tt.name, err)
			}

			// delete
			tmpDataID := list.Items[0].Id // memory id it needs to del later
			delReq := &pb.DelItemRequest{
				DataID: list.Items[0].Id,
			}
			_, err = c.Client.DelItem(ctx, delReq, grpc.Header(&header))
			if err != nil {
				t.Errorf("Error delete information %v %v", tt.name, err)
			}

			// check output after delete
			list, err = c.Client.ListItems(ctx, listReq, grpc.Header(&header))
			if err != nil {
				t.Errorf("Error list information after delete %v %v", tt.name, err)
			}

			if len(list.Items) != 0 {
				t.Errorf("Error length must be 0 after delete %v %v", tt.name, err)
			}

			// clear data
			_ = tearDown.DropData(int(tmpDataID))

			// check stream
			path := "../../cmd/gophkeeper/client/main.go"
			ctx, cancel := context.WithCancel(context.Background())
			err = c.UploadFile(ctx, cancel, 3, path, "", false, 0)
			if err != nil {
				t.Errorf("Error stream upload %v %v", tt.name, err)
			}

			// check output after stream
			ctx = metadata.NewOutgoingContext(context.Background(), md)
			list, err = c.Client.ListItems(ctx, listReq, grpc.Header(&header))
			if err != nil {
				t.Errorf("Error list information after stream %v %v", tt.name, err)
			}

			if len(list.Items) != 1 {
				t.Errorf("Error length must be 1 %v %v", tt.name, err)
			}

			// update stream
			path = "../../cmd/gophkeeper/client/bin/client-darwin-m1"
			ctx, cancel = context.WithCancel(context.Background())
			err = c.UploadFile(ctx, cancel, 3, path, "", false, int(list.Items[0].Id))
			if err != nil {
				t.Errorf("Error stream upload %v %v", tt.name, err)
			}

			// drop test data
			ctx = metadata.NewOutgoingContext(context.Background(), md)
			list, err = c.Client.ListItems(ctx, listReq, grpc.Header(&header))
			if err != nil {
				t.Errorf("Error list before drop data %v %v", tt.name, err)
			}

			for _, v := range list.Items {
				tearDown.DropData(int(v.Id))
			}

			// drop test user
			_ = tearDown.DropUser(resp.User.Login)

		})
	}

}
