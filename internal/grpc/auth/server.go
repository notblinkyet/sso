package authgrpc

import (
	"context"
	"errors"
	"time"

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
}

type serverApi struct {
	ssov1.UnimplementedAuthServer
	auth    Auth
	timeout time.Duration
}

func Register(gRRPCserver *grpc.Server, auth Auth, timeout time.Duration) {
	ssov1.RegisterAuthServer(gRRPCserver, serverApi{auth: auth, timeout: timeout})
}

// Register implements ssov1.AuthServer.
func (s serverApi) Register(ctx context.Context, in *ssov1.RegisterRequest) (*ssov1.RegisterResponse, error) {

	ctx, cancelCtx := context.WithTimeout(ctx, s.timeout)
	defer cancelCtx()

	if in.Login == "" {
		return nil, status.Error(codes.InvalidArgument, "login is required")
	}

	if in.Password == "" {
		return nil, status.Error(codes.InvalidArgument, "password is required")
	}

	uid, err := s.auth.Register(ctx, in.Login, in.Password)

	if err != nil {

		if errors.Is(err, storage.ErrLoginExists) {
			return nil, status.Error(codes.AlreadyExists, "login already exists")
		}

		return nil, status.Error(codes.Internal, "failed to register")
	}

	select {
	case <-ctx.Done():
		return nil, status.Error(codes.DeadlineExceeded, ctx.Err().Error())
	default:
		return &ssov1.RegisterResponse{
			UserId: uid,
		}, nil
	}
}

// Login implements ssov1.AuthServer.
func (s serverApi) Login(ctx context.Context, in *ssov1.LoginRequest) (*ssov1.LoginResponse, error) {

	ctx, cancelCtx := context.WithTimeout(ctx, s.timeout)
	defer cancelCtx()

	token, err := s.auth.Login(ctx, in.Login, in.Password, int(in.AppId))

	if err != nil {
		if errors.Is(err, auth.ErrInvalidCredentials) {
			return nil, status.Error(codes.InvalidArgument, "invalid login or password")
		}
		return nil, status.Error(codes.Internal, "failed to login")
	}

	select {
	case <-ctx.Done():
		return nil, status.Error(codes.DeadlineExceeded, ctx.Err().Error())
	default:
		return &ssov1.LoginResponse{Token: token}, nil
	}
}

// IsAdmin implements ssov1.AuthServer.
func (s serverApi) IsAdmin(ctx context.Context, in *ssov1.IsAdminRequest) (*ssov1.IsAdminResponse, error) {

	ctx, cancelCtx := context.WithTimeout(ctx, s.timeout)
	defer cancelCtx()

	if in.UserId == 0 {
		return nil, status.Error(codes.InvalidArgument, "userId is required")
	}
	isAdmin, err := s.auth.IsAdmin(ctx, in.UserId)

	if err != nil {
		return nil, status.Error(codes.Internal, "failed to get info about rules")
	}

	select {
	case <-ctx.Done():
		return nil, status.Error(codes.DeadlineExceeded, ctx.Err().Error())
	default:
		return &ssov1.IsAdminResponse{IsAdmin: isAdmin}, nil
	}

}
