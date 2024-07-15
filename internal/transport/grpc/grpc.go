package grpc

import (
	"context"
	"net/mail"
	"strings"
	"unicode/utf8"

	pb "github.com/rogue0026/proto/gen/go/sso"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const forbiddenLoginSymbols string = "!@#$%^&*()-_+="

type serverAPI struct {
	pb.UnimplementedAuthServer
	a Auth
}

type Auth interface {
	RegisterNewUser(ctx context.Context, login string, password string, email string) (int64, error)
	LoginUser(ctx context.Context, login string, password string) (string, error)
}

func (s *serverAPI) Register(ctx context.Context, in *pb.RegisterUserRequest) (*pb.RegisterUserResponse, error) {

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
	userId, err := s.a.RegisterNewUser(ctx, in.Login, in.Password, in.Email)
	if err != nil {
		return nil, status.Error(codes.Internal, "internal server error")
	}

	// sending response to client
	resp := pb.RegisterUserResponse{
		UserId: userId,
	}
	return &resp, nil
}

func (s *serverAPI) Login(ctx context.Context, in *pb.LoginUserRequest) (*pb.LoginuserResponse, error) {
	token, err := s.a.LoginUser(ctx, in.Login, in.Password)
	if err != nil {
		// todo
		// 1. Неверный логин или пароль
		// 2. Пользователь не найден => неверный логин или пароль (это усложнит жизнь злоумышленникам, пытающимся получить неправомерный доступ к системе)
		// 3. Внутрення ошибка сервиса
		return nil, status.Error(codes.Internal, "internal server error")
	}

	resp := pb.LoginuserResponse{
		Token: token,
	}

	return &resp, nil
}
