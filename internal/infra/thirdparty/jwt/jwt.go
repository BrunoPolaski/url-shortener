package jwt

import (
	"github.com/BrunoPolaski/go-crud/src/configuration/rest_err"
	"github.com/golang-jwt/jwt/v5"
)

type JWT interface {
	GenerateToken(uid string) (string, *rest_err.RestErr)
	ParseToken(token string) (*jwt.Token, *rest_err.RestErr)
}
