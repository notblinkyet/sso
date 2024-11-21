package suite

import (
	"context"
	"net"
	"os"
	"strconv"
	"testing"

	ssov1 "github.com/notblinkyet/proto_sso/gen/go/sso"
	"github.com/notblinkyet/sso/internal/config"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Suite struct {
	*testing.T
	Cfg        *config.Config
	AuthClient ssov1.AuthClient
}

func configPath() string {

	if v := os.Getenv("CONFIG_PATH"); v != "" {
		return v
	}

	return "/home/hobonail/go_projects/sso_project/sso/config/local.yaml"
}

func New(t *testing.T) (context.Context, *Suite) {
	t.Helper()
	t.Parallel()

	cfg := config.MustLoadFromPath(configPath())

	ctx, cancelCtx := context.WithTimeout(context.Background(), cfg.Grpc.Timeout)

	t.Cleanup(func() {
		t.Helper()
		cancelCtx()
	})

	clienConn, err := grpc.NewClient(net.JoinHostPort(cfg.Grpc.Host, strconv.Itoa(cfg.Grpc.Port)),
		grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		t.Fatalf("grpc server connection failed: %v", err)
	}

	return ctx, &Suite{
		T:          t,
		Cfg:        cfg,
		AuthClient: ssov1.NewAuthClient(clienConn),
	}

}
