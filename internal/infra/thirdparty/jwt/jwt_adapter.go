package jwt

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/BrunoPolaski/go-crud/src/configuration/rest_err"
	"github.com/golang-jwt/jwt/v5"
)

type jwtAdapter struct{}

func NewJWTAdapter() JWT {
	return &jwtAdapter{}
}

func (ja *jwtAdapter) GenerateToken(sub string) (string, *rest_err.RestErr) {
	secret := os.Getenv("TOKEN_SECRET")
	if secret == "" {
		return "", rest_err.NewInternalServerError("token secret is not set")
	}

	expTime, err := strconv.Atoi(os.Getenv("JWT_EXPIRATION_TIME"))
	if err != nil {
		expTime = 3600
	}

	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		jwt.MapClaims{
			"sub": sub,
			"iss": os.Getenv("APP_URL"),
			"exp": time.Now().Add(time.Second * time.Duration(expTime)).Unix(),
			"iat": time.Now().Unix(),
		},
	)

	tokenWithSignature, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", rest_err.NewInternalServerError(fmt.Sprintf("Failed to sign JWT: %s", err.Error()))
	}

	return tokenWithSignature, nil
}

func (ja *jwtAdapter) ParseToken(token string) (*jwt.Token, *rest_err.RestErr) {
	secret := os.Getenv("TOKEN_SECRET")

	parsedToken, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); ok {
			return []byte(secret), nil
		}

		return nil, rest_err.NewBadRequestError("invalid token")
	})

	if err != nil {
		return nil, rest_err.NewBadRequestError("invalid token")
	}

	return parsedToken, nil
}
