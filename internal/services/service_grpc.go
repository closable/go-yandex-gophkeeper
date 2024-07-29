package services

import (
	"context"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/closable/go-yandex-gophkeeper/internal/errors"
	pb "github.com/closable/go-yandex-gophkeeper/internal/services/proto"
	"github.com/closable/go-yandex-gophkeeper/internal/store"
	"github.com/closable/go-yandex-gophkeeper/internal/utils"
	"go.uber.org/zap"
	"google.golang.org/genproto/googleapis/rpc/code"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type (
	// GRPCStorager main service interface
	GRPCStorager interface {
		//PrepareDB() error
		AddItem(userId, dataType int, data, name string) error
		GetUserInfo(login, password string) (*store.UserDetail, error)
		ListItems(userId int) ([]store.RowItem, error)
		UpdateItem(userId, dataId int, data string) error
		DeleteItem(userId, dataId int) error

		GetUserKeyString(userID int) (string, error)
		CreateUser(user, pass, keyStr string) (*store.UserDetail, error)
		CheckUser(user string) bool
		Health(n string) error
	}

	// GRPCFileStorager main fileservice interface
	GRPCFileStorager interface {
		Upload(stream pb.FilseService_UploadServer) (*pb.FileUploadResponse, error)
		Download(in *pb.FileDownloadRequest, srv pb.FilseService_DownloadServer) error
		AddItem(userId, dataType int, data, name string) error
		UpdateItem(userId, dataId int, data string) error
		GetUserKeyString(userID int) (string, error)
		GetFileData(dataID int) (*store.FileData, error)
	}

	// GophKeeperServer main server structure
	GophKeeperServer struct {
		pb.UnimplementedGophKeeperServer
		store  GRPCStorager
		addr   string
		logger *zap.Logger
	}

	// GophKeeperFileServer main fie server structure
	GophKeeperFileServer struct {
		pb.UnimplementedFilseServiceServer
		store  GRPCFileStorager
		addr   string
		logger *zap.Logger
	}
)

// Health server status
func (s *GophKeeperServer) Health(ctx context.Context, in *pb.HealthRequest) (*pb.HealthResponse, error) {
	var response pb.HealthResponse

	err := s.store.Health(in.Numb)
	if err != nil {
		s.logger.Info(fmt.Sprintf("Server status %s", err))
		return &response, err
	}
	s.logger.Info(fmt.Sprintf("Server status %s", "online"))
	return &response, nil
}

// Login user auth
func (s *GophKeeperServer) Login(ctx context.Context, in *pb.LoginRequest) (*pb.LoginResponse, error) {
	var response pb.LoginResponse

	user, err := s.store.GetUserInfo(in.User, in.Pass)
	if err != nil {
		return &response, err
	}

	token, err := utils.BuildJWTString(user.UserID)
	if err != nil {
		s.logger.Info(fmt.Sprintf("%v %s", errors.ErrorJWTToken, err))
		return &response, status.Errorf(codes.Code(code.Code_INTERNAL), fmt.Sprintf("%v %s", errors.ErrorJWTToken, err))
	}
	s.logger.Info(fmt.Sprintf("JWT token created %s", "OK"))
	response.Token = token
	return &response, nil
}

// CreateUser create new user
func (s *GophKeeperServer) CreateUser(ctx context.Context, in *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	var response pb.CreateUserResponse

	if len(in.User) == 0 || len(in.Pass) == 0 {
		s.logger.Info(fmt.Sprintf("%v %s %s", errors.ErrorAuthInfo, in.User, in.Pass))
		return &response, status.Errorf(codes.Code(code.Code_INVALID_ARGUMENT), errors.ErrorAuthInfo.Error())
	}
	keyString := ""
	keyString = in.Keystring
	if len(in.Keystring) == 0 {
		keyString = utils.GetRandomString(32)
	}
	key, err := utils.CryptoSeq(keyString)
	if err != nil {
		s.logger.Info(fmt.Sprintf("%v %s %s", errors.ErrorCrypoSeq, keyString, err))
		return &response, status.Errorf(codes.Code(code.Code_INTERNAL), err.Error())
	}

	usr, err := s.store.CreateUser(in.User, in.Pass, key)
	if err != nil {
		s.logger.Info(fmt.Sprintf("%v %s %s %v", errors.ErrorAddItem, in.User, in.Pass, err))
		return &response, status.Errorf(codes.Code(code.Code_INTERNAL), fmt.Sprintf("%v %s %s %v", errors.ErrorAddItem, in.User, in.Pass, err))
	}

	token, err := utils.BuildJWTString(usr.UserID)
	if err != nil {
		s.logger.Info(fmt.Sprintf("%sv %v", errors.ErrorJWTToken.Error(), err))
		return &response, status.Errorf(codes.Code(code.Code_INTERNAL), err.Error())
	}

	s.logger.Info(fmt.Sprintf("Create user %v", "OK"))
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
		s.logger.Error(fmt.Sprintf("%v %v %v", errors.ErrorListData, userID, err))
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
	s.logger.Info(fmt.Sprintf("List items %v", "OK"))
	return &response, nil
}

// DeleteItem delete item
func (s *GophKeeperServer) DelItem(ctx context.Context, in *pb.DelItemRequest) (*pb.DelItemResponse, error) {
	response := &pb.DelItemResponse{}

	userID := ctx.Value("user_id").(int)

	err := s.store.DeleteItem(userID, int(in.DataID))
	if err != nil {
		s.logger.Error(fmt.Sprintf("%v %v %v", errors.ErrorDeleteData, userID, err))
		return response, err
	}
	s.logger.Info(fmt.Sprintf("Item deleted %v", "OK"))
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
		s.logger.Error(fmt.Sprintf("%v %v %v", errors.ErrorCrypoSeq, userID, err))
		return &response, err
	}

	encData := utils.Encrypt(keyString, in.Data)

	err = s.store.AddItem(userID, int(in.DataType), encData, in.Name)
	if err != nil {
		s.logger.Error(fmt.Sprintf("%v %v %v", errors.ErrorAddItem, userID, err))
		return &response, err
	}
	s.logger.Info(fmt.Sprintf("Item added %v %v", userID, "OK"))
	return &response, nil
}

// Upload file to server
func (s *GophKeeperFileServer) Upload(stream pb.FilseService_UploadServer) error {
	fileSize := uint32(0)

	file, err := os.CreateTemp("", fmt.Sprintf("tmp_%s", utils.GetRandomString(10)))
	if err != nil {
		s.logger.Error(fmt.Sprintf("%v %v", errors.ErrorFileCreate, err))
		return err
	}
	defer file.Close()
	defer os.Remove(file.Name())
	dataID := int32(0)
	addItemReq := &pb.AddItemWithTokenRequest{}
	flag := false
	s.logger.Info(fmt.Sprintf("Getting stream data .... %s", ""))
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
			s.logger.Error(fmt.Sprintf("%v %v", errors.ErrorStreamData, err))
			return status.Errorf(codes.Code(code.Code_INTERNAL), err.Error())
		}
		chunk := req.GetChunk()
		fileSize += uint32(len(chunk))
		_, err = file.Write(chunk)
		if err != nil {
			s.logger.Error(fmt.Sprintf("%v %v", errors.ErrorFileWriting, err))
			return status.Errorf(codes.Code(code.Code_INTERNAL), err.Error())
		}
	}

	bin := make([]byte, fileSize)
	_, err = file.ReadAt(bin, 0)
	if err != nil {
		s.logger.Error(fmt.Sprintf("%v %v", errors.ErrorFileReading, err))
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
		s.logger.Error(fmt.Sprintf("%v %v", errors.ErrorFileWriting, err))
		return status.Errorf(codes.Code(code.Code_INTERNAL), err.Error())
	}
	s.logger.Info(fmt.Sprintf("Operation new(updated) %v", "OK"))
	return nil
}

// AddItem add new item into store
func (s *GophKeeperFileServer) AddItem(ctx context.Context, in *pb.AddItemWithTokenRequest) (*pb.AddItemResponse, error) {
	var response pb.AddItemResponse

	userID := utils.GetUserID(in.Token)
	if userID == 0 {
		s.logger.Error(fmt.Sprintf("%v %v", errors.ErrorJWTToken.Error(), in.Token))
		return &response, status.Errorf(codes.Unauthenticated, "%v: %v", errors.ErrorJWTToken.Error(), in.Token)
	}

	keyString, err := s.store.GetUserKeyString(userID)
	if err != nil {
		s.logger.Error(fmt.Sprintf("%v %v", errors.ErrorCrypoSeq, userID))
		return &response, err
	}

	encData := utils.Encrypt(keyString, in.Data)
	err = s.store.AddItem(userID, int(in.DataType), encData, in.Name)
	if err != nil {
		s.logger.Error(fmt.Sprintf("%v %v %v", errors.ErrorAddItem, userID, err))
		return &response, err
	}
	s.logger.Info(fmt.Sprintf("Item added %v %v", userID, "OK"))
	return &response, nil
}

// UpdateItem update item
func (s *GophKeeperFileServer) UpdateItem(ctx context.Context, in *pb.UpdateItemWithTokenRequest) (*pb.UpdateItemResponse, error) {
	var response pb.UpdateItemResponse

	userID := utils.GetUserID(in.Token)
	if userID == 0 {
		s.logger.Error(fmt.Sprintf("%v %v", errors.ErrorJWTToken.Error(), in.Token))
		return &response, status.Errorf(codes.Unauthenticated, "%v: %v", errors.ErrorJWTToken.Error(), in.Token)
	}

	key, err := s.store.GetUserKeyString(userID)
	if err != nil {
		s.logger.Error(fmt.Sprintf("%v %v", errors.ErrorCrypoSeq, userID))
		return &response, err
	}

	encData := utils.Encrypt(key, in.Data)

	err = s.store.UpdateItem(userID, int(in.DataID), encData)
	if err != nil {
		s.logger.Error(fmt.Sprintf("%v %v %v", errors.ErrorUpdateItem, userID, err))
		return &response, err
	}
	s.logger.Info(fmt.Sprintf("Item updated %v %v", userID, "OK"))
	return &response, nil
}

func (s *GophKeeperFileServer) Download(in *pb.FileDownloadRequest, srv pb.FilseService_DownloadServer) error {
	userID := utils.GetUserID(in.Token)
	key, err := s.store.GetUserKeyString(userID)
	if err != nil {
		return status.Errorf(codes.Code(code.Code_INTERNAL), err.Error())
	}
	data, err := s.store.GetFileData(int(in.DataID))
	if err != nil || len(data.Data) == 0 {
		return status.Errorf(codes.Code(code.Code_INTERNAL), err.Error())
	}

	reader := strings.NewReader(utils.Decrypt(key, data.Data))
	buf := make([]byte, 1024)
	batchNumber := 1

	for {
		num, err := reader.Read(buf)
		if err == io.EOF {
			break
		}
		if err != nil {
			return status.Errorf(codes.Code(code.Code_INTERNAL), err.Error())
		}
		chunk := buf[:num]
		err = srv.Send(&pb.FileDownloadResponse{FilePath: data.FilePath, DataType: int32(data.DataType), Chank: chunk})
		if err != nil {
			return status.Errorf(codes.Code(code.Code_INTERNAL), err.Error())
		}
		batchNumber += 1
	}
	s.logger.Info(fmt.Sprintf("File downloaded %v %v", data.FilePath, "OK"))
	return nil
}
