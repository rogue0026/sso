package grpc

import (
	"context"
	"errors"
	"log/slog"
	"net/mail"
	"strings"
	"unicode/utf8"

	pb "github.com/rogue0026/proto/gen/go/sso"
	"github.com/rogue0026/sso/internal/services/auth"
	"github.com/rogue0026/sso/internal/storage"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const forbiddenLoginSymbols string = `"!@#$%^&*()-_+=`

type grpcAPI struct {
	Logger *slog.Logger
	pb.UnimplementedAuthServer
	auth Auth
}

func NewAPI(logger *slog.Logger, service Auth) *grpcAPI {
	api := grpcAPI{
		Logger: logger,
		auth:   service,
	}
	return &api
}

type Auth interface {
	RegisterNewUser(ctx context.Context, login string, password string, email string) (int64, error)
	LoginUser(ctx context.Context, login string, password string) (string, error)
}

func (s *grpcAPI) Register(ctx context.Context, in *pb.RegisterUserRequest) (*pb.RegisterUserResponse, error) {

	// make validation input params
	_, err := mail.ParseAddress(in.Email)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid email address")
	}

	passLen := utf8.RuneCountInString(in.Password)
	if passLen < 8 {
		return nil, status.Error(codes.InvalidArgument, "password length must be equal or greater than 8 symbols")
	}

	if strings.ContainsAny(in.Login, forbiddenLoginSymbols) {
		return nil, status.Error(codes.InvalidArgument, "login contains forbiden symbols")
	}

	// if validation ok, send call to service layer
	userId, err := s.auth.RegisterNewUser(ctx, in.Login, in.Password, in.Email)
	if err != nil {
		return nil, status.Error(codes.Internal, "internal server error")
	}

	// sending response to client
	resp := pb.RegisterUserResponse{
		UserId: userId,
	}
	return &resp, nil
}

func (s *grpcAPI) Login(ctx context.Context, in *pb.LoginUserRequest) (*pb.LoginuserResponse, error) {
	token, err := s.auth.LoginUser(ctx, in.Login, in.Password)
	if err != nil {
		if errors.Is(err, auth.ErrInvalidUserCredentials) {
			return nil, status.Error(codes.Unauthenticated, auth.ErrInvalidUserCredentials.Error())
		}
		if errors.Is(err, storage.ErrUserNotFound) {
			return nil, status.Error(codes.NotFound, storage.ErrUserNotFound.Error())
		}
		return nil, status.Error(codes.Internal, "internal server error")
	}

	resp := pb.LoginuserResponse{
		Token: token,
	}

	return &resp, nil
}
