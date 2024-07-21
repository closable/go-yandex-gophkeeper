package services

import (
	"context"
	"fmt"
	"io"
	"os"

	pb "github.com/closable/go-yandex-gophkeeper/internal/services/proto"
	"github.com/closable/go-yandex-gophkeeper/internal/store"
	"github.com/closable/go-yandex-gophkeeper/internal/utils"
	"google.golang.org/genproto/googleapis/rpc/code"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type (
	GRPCStorager interface {
		AddItem(userId, dataType int, data, name string) error
		GetUserInfo(login, password string) (*store.UserDetail, error)
		ListItems(userId int) ([]store.RowItem, error)
		UpdateItem(userId, dataId int, data string) error
		DeleteItem(userId, dataId int) error

		GetUserKeyString(userID int) (string, error)
		CreateUser(user, pass, keyStr string) (*store.UserDetail, error)
		CheckUser(user string) bool
	}

	GRPCFileStorager interface {
		Upload(stream pb.FilseService_UploadServer) (*pb.FileUploadResponse, error)
		AddItem(userId, dataType int, data, name string) error
		UpdateItem(userId, dataId int, data string) error
		GetUserKeyString(userID int) (string, error)
	}

	GophKeeperServer struct {
		pb.UnimplementedGophKeeperServer
		store GRPCStorager
		addr  string
	}

	GophKeeperFileServer struct {
		pb.UnimplementedFilseServiceServer
		store GRPCFileStorager
		addr  string
	}
)

// Login user auth
func (s *GophKeeperServer) Login(ctx context.Context, in *pb.LoginRequest) (*pb.LoginResponse, error) {
	var response pb.LoginResponse

	user, err := s.store.GetUserInfo(in.User, in.Pass)
	if err != nil {
		return &response, err
	}

	token, err := utils.BuildJWTString(user.UserID)
	if err != nil {
		return &response, status.Errorf(codes.Code(code.Code_INTERNAL), "error")
	}

	response.Token = token
	return &response, nil
}

// CreateUser create new user
func (s *GophKeeperServer) CreateUser(ctx context.Context, in *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	var response pb.CreateUserResponse

	if len(in.User) == 0 || len(in.Pass) == 0 {
		return &response, status.Errorf(codes.Code(code.Code_INVALID_ARGUMENT), "login or password empty")
	}
	keyString := ""
	keyString = in.Keystring
	if len(in.Keystring) == 0 {
		keyString = utils.GetRandomString(32)
	}
	key, err := utils.CryptoSeq(keyString)
	if err != nil {
		return &response, status.Errorf(codes.Code(code.Code_INTERNAL), "error")
	}

	usr, err := s.store.CreateUser(in.User, in.Pass, key)
	if err != nil {
		return &response, status.Errorf(codes.Code(code.Code_INTERNAL), "internal server error")
	}

	token, err := utils.BuildJWTString(usr.UserID)
	if err != nil {
		return &response, status.Errorf(codes.Code(code.Code_INTERNAL), "error")
	}

	return &pb.CreateUserResponse{
		User: &pb.UserDetail{
			UserID:    int32(usr.UserID),
			Login:     usr.Login,
			Keystring: key,
		},
		Token: token,
	}, nil
}

// ListItems list items
func (s *GophKeeperServer) ListItems(ctx context.Context, in *pb.ListItemsRequest) (*pb.ListItemsResponse, error) {
	var response pb.ListItemsResponse
	var data []*pb.Item

	userID := ctx.Value("user_id").(int)

	key, err := s.store.GetUserKeyString(userID)
	if err != nil {
		return &response, err
	}

	rows, err := s.store.ListItems(userID)
	if err != nil {
		return &response, err
	}

	for _, v := range rows {
		txt := v.EncData
		if in.Decrypted && v.DataType < 3 {
			txt = utils.Decrypt(key, v.EncData)
		}
		data = append(data, &pb.Item{
			Id:       int32(v.Id),
			Type:     v.Type,
			Name:     v.Name,
			Restore:  v.IsRestore,
			Encdata:  txt,
			Length:   int32(v.Length),
			DataType: int32(v.DataType),
		})
	}
	response.Items = data
	return &response, nil
}

// DeleteItem delete item
func (s *GophKeeperServer) DelItem(ctx context.Context, in *pb.DelItemRequest) (*pb.DelItemResponse, error) {
	response := &pb.DelItemResponse{}

	userID := ctx.Value("user_id").(int)

	err := s.store.DeleteItem(userID, int(in.DataID))
	if err != nil {
		return response, err
	}
	return response, nil
}

// UpdateItem update item
func (s *GophKeeperServer) UpdateItem(ctx context.Context, in *pb.UpdateItemRequest) (*pb.UpdateItemResponse, error) {
	var response pb.UpdateItemResponse

	userID := ctx.Value("user_id").(int)

	key, err := s.store.GetUserKeyString(userID)
	if err != nil {
		return &response, err
	}

	encData := utils.Encrypt(key, in.Data)

	err = s.store.UpdateItem(userID, int(in.DataID), encData)
	if err != nil {
		return &response, err
	}

	return &response, nil
}

// AddItem add new item into store
func (s *GophKeeperServer) AddItem(ctx context.Context, in *pb.AddItemRequest) (*pb.AddItemResponse, error) {
	var response pb.AddItemResponse

	userID := ctx.Value("user_id").(int)

	keyString, err := s.store.GetUserKeyString(userID)
	if err != nil {
		return &response, err
	}

	encData := utils.Encrypt(keyString, in.Data)

	err = s.store.AddItem(userID, int(in.DataType), encData, in.Name)
	if err != nil {
		return &response, err
	}

	return &response, nil
}

func (s *GophKeeperFileServer) Upload(stream pb.FilseService_UploadServer) error {
	fileSize := uint32(0)

	file, err := os.CreateTemp("", fmt.Sprintf("tmp_%s", utils.GetRandomString(10)))
	if err != nil {
		return err
	}
	defer file.Close()
	defer os.Remove(file.Name())
	dataID := int32(0)
	addItemReq := &pb.AddItemWithTokenRequest{}
	flag := false

	for {
		req, err := stream.Recv()

		if !flag {
			addItemReq.Token = req.GetToken()
			addItemReq.DataType = req.GetDataType()
			addItemReq.Data = req.GetData()
			addItemReq.Name = req.GetName()
			dataID = req.GetDataID() // 0 - new int update
			flag = true
		}

		if err == io.EOF {
			break
		}
		if err != nil {
			return status.Errorf(codes.Code(code.Code_INTERNAL), err.Error())
		}
		chunk := req.GetChunk()
		fileSize += uint32(len(chunk))
		_, err = file.Write(chunk)
		if err != nil {
			return status.Errorf(codes.Code(code.Code_INTERNAL), err.Error())
		}
	}

	bin := make([]byte, fileSize)
	_, err = file.ReadAt(bin, 0)
	if err != nil {
		return status.Errorf(codes.Code(code.Code_INTERNAL), err.Error())
	}
	addItemReq.Data = string(bin)

	if dataID == 0 {
		_, err = s.AddItem(context.Background(), addItemReq)
	} else {
		updItemReques := &pb.UpdateItemWithTokenRequest{
			Token:  addItemReq.Token,
			DataID: dataID,
			Data:   addItemReq.Data,
		}
		_, err = s.UpdateItem(context.Background(), updItemReques)
	}

	if err != nil {
		return status.Errorf(codes.Code(code.Code_INTERNAL), err.Error())
	}
	return nil
}

// AddItem add new item into store
func (s *GophKeeperFileServer) AddItem(ctx context.Context, in *pb.AddItemWithTokenRequest) (*pb.AddItemResponse, error) {
	var response pb.AddItemResponse

	userID := utils.GetUserID(in.Token)
	if userID == 0 {
		return &response, status.Errorf(codes.Unauthenticated, "invalid auth token: %v", in.Token)
	}

	keyString, err := s.store.GetUserKeyString(userID)
	if err != nil {
		return &response, err
	}

	encData := utils.Encrypt(keyString, in.Data)
	err = s.store.AddItem(userID, int(in.DataType), encData, in.Name)
	if err != nil {
		return &response, err
	}

	return &response, nil
}

// UpdateItem update item
func (s *GophKeeperFileServer) UpdateItem(ctx context.Context, in *pb.UpdateItemWithTokenRequest) (*pb.UpdateItemResponse, error) {
	var response pb.UpdateItemResponse

	userID := utils.GetUserID(in.Token)
	if userID == 0 {
		return &response, status.Errorf(codes.Unauthenticated, "invalid auth token: %v", in.Token)
	}

	key, err := s.store.GetUserKeyString(userID)
	if err != nil {
		return &response, err
	}

	encData := utils.Encrypt(key, in.Data)

	err = s.store.UpdateItem(userID, int(in.DataID), encData)
	if err != nil {
		return &response, err
	}

	return &response, nil
}
