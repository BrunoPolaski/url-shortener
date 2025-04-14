package repositories

import "github.com/BrunoPolaski/go-crud/src/configuration/rest_err"

type AuthRepository interface {
	FindToken(uuid string) (string, *rest_err.RestErr)
}
