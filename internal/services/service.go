package services

import (
	"context"
	"fmt"
	"net"
	"strings"

	"github.com/closable/go-yandex-gophkeeper/internal/store"
	"github.com/closable/go-yandex-gophkeeper/internal/utils"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/closable/go-yandex-gophkeeper/internal/services/proto"
	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
)

// ServTransporter server structire
type ServTransporter struct {
	ServAddr         string
	FileServAddr     string
	store            GRPCStorager
	GRPCServ         *grpc.Server
	GRPCFileServ     *grpc.Server
	GRPCFileListener net.Listener
	GRPCListener     net.Listener
	Logger           *zap.Logger
}

// New new instance server
func New(DSN, addr, addrFileServ string, logger *zap.Logger) (*ServTransporter, error) {
	var st GRPCStorager
	var fileSt GRPCFileStorager
	st, err := store.New(DSN)
	if err != nil {
		logger.Error(err.Error())
		panic(err)
	}

	fileSt, err = store.New(DSN)
	if err != nil {
		logger.Error(err.Error())
		panic(err)
	}
	return GRPCserv(st, fileSt, addr, addrFileServ, logger)
}

// GRPCserv configure grpc server
func GRPCserv(store GRPCStorager, fileStore GRPCFileStorager, addr, addrFileServ string, logger *zap.Logger) (*ServTransporter, error) {
	listen, err := net.Listen("tcp", addr)
	if err != nil {
		logger.Error(err.Error())
		panic(err)
	}

	fileListen, err := net.Listen("tcp", addrFileServ)
	if err != nil {
		logger.Error(err.Error())
		panic(err)
	}

	serv := grpc.NewServer(
		grpc.StreamInterceptor(grpc_auth.StreamServerInterceptor(AuthFunc)),
		grpc.UnaryInterceptor(UnaryServerInterceptor(AuthFunc)),
	)
	// регистрируем сервис
	pb.RegisterGophKeeperServer(serv, &GophKeeperServer{
		store:  store,
		addr:   addr,
		logger: logger,
	})

	fileServ := grpc.NewServer()
	// регистрируем файловый сервис
	pb.RegisterFilseServiceServer(fileServ, &GophKeeperFileServer{
		store:  fileStore,
		addr:   addr,
		logger: logger,
	})

	logger.Info("Server structure is ready!")
	return &ServTransporter{
		ServAddr:         addr,
		FileServAddr:     addrFileServ,
		store:            store,
		GRPCServ:         serv,
		GRPCFileServ:     fileServ,
		GRPCFileListener: fileListen,
		GRPCListener:     listen,
		Logger:           logger,
	}, nil

}

// Run start server service
func (t *ServTransporter) Run() error {

	go func() {
		err := t.GRPCFileServ.Serve(t.GRPCFileListener)
		if err != nil {
			//fmt.Printf("File Server error %s", err)
			t.Logger.Info(fmt.Sprintf("File Server error %s", err))
			panic(err)
		}
		//fmt.Printf("File Server started %s", t.FileServAddr)
		t.Logger.Info(fmt.Sprintf("File Server started %s", t.FileServAddr))
	}()

	fmt.Println("Server started ", t.ServAddr)
	t.Logger.Info(fmt.Sprintf("Server started %s", t.ServAddr))
	return t.GRPCServ.Serve(t.GRPCListener)
}

// Shutdown use for Graceful shutdown
func (t *ServTransporter) Shutdown(ctx context.Context) error {
	t.GRPCServ.GracefulStop()
	t.GRPCFileServ.GracefulStop()
	return nil
}

// AuthFunc check auth function (middleware)
func AuthFunc(ctx context.Context) (context.Context, error) {
	token, err := grpc_auth.AuthFromMD(ctx, "bearer")
	if err != nil {
		return nil, err
	}

	userID := utils.GetUserID(token)
	if userID == 0 {
		return nil, status.Errorf(codes.Unauthenticated, "invalid auth token: %v", token)
	}

	grpc_ctxtags.Extract(ctx).Set("auth.sub", userID)
	newCtx := context.WithValue(ctx, "user_id", userID)

	return newCtx, nil
}

// UnaryServerInterceptor
func UnaryServerInterceptor(authFunc grpc_auth.AuthFunc) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		var newCtx context.Context
		var err error

		m := strings.Split(info.FullMethod, "/")
		method := m[len(m)-1]

		if method == "Login" || method == "CreateUser" {
			newCtx = ctx
		} else {
			newCtx, err = authFunc(ctx)
		}

		if err != nil {
			return nil, err
		}
		return handler(newCtx, req)
	}
}
