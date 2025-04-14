package services

import "github.com/BrunoPolaski/go-crud/src/configuration/rest_err"

type AuthService interface {
	RefreshToken() (string, *rest_err.RestErr)
}
