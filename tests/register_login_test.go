package tests

import (
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/golang-jwt/jwt/v5"
	ssov1 "github.com/notblinkyet/proto_sso/gen/go/sso"
	"github.com/notblinkyet/sso/tests/suite"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	DefaultLengthPass int   = 10
	goodAppid         int32 = 1
	wrongAppid        int32 = 0
	appSecret               = "test"
)

func passfunc(length int) string {
	return gofakeit.Password(true, true, true, true, false, length)
}

func TestHappy_Register_Login(t *testing.T) {
	ctx, suite := suite.New(t)

	login := gofakeit.Username()
	password := passfunc(DefaultLengthPass)

	regResp, err := suite.AuthClient.Register(ctx, &ssov1.RegisterRequest{
		Login:    login,
		Password: password,
	})

	require.NoError(t, err)
	assert.NotEmpty(t, regResp.GetUserId())

	logResp, err := suite.AuthClient.Login(ctx, &ssov1.LoginRequest{
		Login:    login,
		Password: password,
		AppId:    goodAppid,
	})

	token := logResp.Token

	require.NoError(t, err)
	assert.NotEmpty(t, token)

	parsedToken, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		return []byte(appSecret), nil
	})

	require.NoError(t, err)

	loginTime := time.Now()

	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	require.True(t, ok)

	assert.Equal(t, regResp.GetUserId(), int64(claims["uid"].(float64)))
	assert.Equal(t, login, claims["login"].(string))
	assert.Equal(t, goodAppid, int32(claims["appID"].(float64)))
	assert.InDelta(t, loginTime.Add(suite.Cfg.TokenTTL).Unix(), claims["exp"].(float64), 1)

}

func TestDublicateLogin(t *testing.T) {
	ctx, suite := suite.New(t)

	login := gofakeit.Username()
	password := passfunc(DefaultLengthPass)

	_, err := suite.AuthClient.Register(ctx, &ssov1.RegisterRequest{
		Login:    login,
		Password: password,
	})

	require.NoError(t, err)

	regResp, err := suite.AuthClient.Register(ctx, &ssov1.RegisterRequest{
		Login:    login,
		Password: password,
	})

	require.Error(t, err)
	assert.Empty(t, regResp.GetUserId())
	assert.ErrorContains(t, err, "login already exists")
}

func TestEmptyLoginPassword(t *testing.T) {
	ctx, suite := suite.New(t)

	login := ""
	password := passfunc(DefaultLengthPass)

	_, err := suite.AuthClient.Register(ctx, &ssov1.RegisterRequest{
		Login:    login,
		Password: password,
	})

	require.Error(t, err)
	assert.ErrorContains(t, err, "login is required")

	login = gofakeit.Username()
	password = ""

	_, err = suite.AuthClient.Register(ctx, &ssov1.RegisterRequest{
		Login:    login,
		Password: password,
	})

	require.Error(t, err)
	assert.ErrorContains(t, err, "password is required")

	login = ""
	password = ""

	_, err = suite.AuthClient.Register(ctx, &ssov1.RegisterRequest{
		Login:    login,
		Password: password,
	})

	require.Error(t, err)
	assert.ErrorContains(t, err, "login is required")
}

func TestWrongApp(t *testing.T) {
	ctx, suite := suite.New(t)

	login := gofakeit.Username()
	password := passfunc(DefaultLengthPass)

	_, err := suite.AuthClient.Register(ctx, &ssov1.RegisterRequest{
		Login:    login,
		Password: password,
	})

	require.NoError(t, err)

	_, err = suite.AuthClient.Login(ctx, &ssov1.LoginRequest{
		Login:    login,
		Password: password,
		AppId:    wrongAppid,
	})
	require.Error(t, err)
}
