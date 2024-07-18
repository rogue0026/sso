package app

import (
	"fmt"
	"log/slog"
	"net"

	pb "github.com/rogue0026/proto/gen/go/sso"
	"github.com/rogue0026/sso/internal/config"
	"github.com/rogue0026/sso/internal/services/auth"
	"github.com/rogue0026/sso/internal/storage/sqlite"
	api "github.com/rogue0026/sso/internal/transport/grpc"
	"google.golang.org/grpc"
)

type Application struct {
	grpcServer *grpc.Server
}

func New(cfg config.Config, logger *slog.Logger) (*Application, error) {

	// init users storage
	users, err := sqlite.New(logger, cfg.DSN)
	if err != nil {
		return nil, err
	}

	// init service object
	authService := auth.New(logger, users, users)

	// creating new grpc server
	grpcServer := grpc.NewServer()

	// creating new api
	grpcAPI := api.NewAPI(logger, authService)

	// register api on grpc server
	pb.RegisterAuthServer(grpcServer, grpcAPI)

	app := Application{
		grpcServer: grpcServer,
	}

	return &app, nil
}

func (a *Application) MustRun(portNumber int) {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", portNumber))
	if err != nil {
		panic(err.Error())
	}

	err = a.grpcServer.Serve(listener)
	if err != nil {
		panic(err.Error())
	}

}
