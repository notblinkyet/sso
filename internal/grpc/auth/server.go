package authgrpc

import (
	"context"
	"errors"

	ssov1 "github.com/notblinkyet/proto_sso/gen/go/sso"
	"github.com/notblinkyet/sso/internal/services/auth"
	storage "github.com/notblinkyet/sso/internal/storage/main_storage"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Auth interface {
	Register(ctx context.Context, login, password string) (int64, error)
	Login(ctx context.Context, login string, password string, appID int) (string, error)
	IsAdmin(ctx context.Context, userID int64) (bool, error)
	Logout(ctx context.Context, token string) (bool, error)
}

type serverApi struct {
	ssov1.UnimplementedAuthServer
	auth Auth
}

func Register(gRRPCserver *grpc.Server, auth Auth) {
	ssov1.RegisterAuthServer(gRRPCserver, serverApi{auth: auth})
}

// IsAdmin implements ssov1.AuthServer.
func (s serverApi) IsAdmin(ctx context.Context, in *ssov1.IsAdminRequest) (*ssov1.IsAdminResponse, error) {

	if in.UserId == 0 {
		return nil, status.Error(codes.InvalidArgument, "userId is required")
	}
	isAdmin, err := s.auth.IsAdmin(ctx, in.UserId)

	if err != nil {
		return nil, status.Error(codes.Internal, "failed to get info about rules")
	}

	return &ssov1.IsAdminResponse{IsAdmin: isAdmin}, nil
}

// Login implements ssov1.AuthServer.
func (s serverApi) Login(ctx context.Context, in *ssov1.LoginRequest) (*ssov1.LoginResponse, error) {
	token, err := s.auth.Login(ctx, in.Login, in.Password, int(in.AppId))

	if err != nil {
		if errors.Is(err, auth.ErrInvalidCredentials) {
			return nil, status.Error(codes.InvalidArgument, "invalid login or password")
		}
		return nil, status.Error(codes.Internal, "failed to login")
	}

	return &ssov1.LoginResponse{Token: token}, nil
}

// Logout implements ssov1.AuthServer.
func (s serverApi) Logout(ctx context.Context, in *ssov1.LogoutRequest) (*ssov1.LogoutResponse, error) {
	success, err := s.auth.Logout(ctx, in.Token)

	if err != nil {
		return nil, status.Error(codes.Internal, "failed to logout")
	}

	return &ssov1.LogoutResponse{Success: success}, nil
}

// Register implements ssov1.AuthServer.
func (s serverApi) Register(ctx context.Context, in *ssov1.RegisterRequest) (*ssov1.RegisterResponse, error) {

	if in.Login == "" {
		return nil, status.Error(codes.InvalidArgument, "login is required")
	}

	if in.Password == "" {
		return nil, status.Error(codes.InvalidArgument, "password is required")
	}

	uid, err := s.auth.Register(ctx, in.Login, in.Password)

	if err != nil {

		if errors.Is(err, storage.ErrUserExists) {
			return nil, status.Error(codes.AlreadyExists, "login already exists")
		}

		return nil, status.Error(codes.Internal, "failed to register")
	}

	return &ssov1.RegisterResponse{
		UserId: uid,
	}, nil
}
