package jwt

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/notblinkyet/sso/internal/models"
)

func NewToken(user *models.User, app *models.App, duration time.Duration) (string, error) {

	token := jwt.New(jwt.SigningMethodHS256)

	claim := token.Claims.(jwt.MapClaims)

	claim["uid"] = user.ID
	claim["login"] = user.Login
	claim["exp"] = time.Now().Add(duration).Unix()
	claim["appID"] = app.ID

	tockenString, err := token.SignedString([]byte(app.Secret))

	if err != nil {

		return "", err
	}

	return tockenString, nil
}
