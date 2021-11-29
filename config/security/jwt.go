package security

import (
	"github.com/dgrijalva/jwt-go"
	"go-module/model"
	"time"
)

const JWT_KEY = "tiger"

func GenToken(user model.User) (string, error) {
	claims := &model.JwtCustomClaims{
		UserId: user.UserId,
		Role:   user.Role,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 1).Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// ky token
	result, err := token.SignedString([]byte(JWT_KEY))
	if err != nil {
		return "", err
	}

	return result, nil
}
