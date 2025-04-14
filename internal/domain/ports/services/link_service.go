package services

import "github.com/BrunoPolaski/go-crud/src/configuration/rest_err"

type LinkService interface {
	Redirect(uuid string) (string, *rest_err.RestErr)
	CreateLink(url string) (string, *rest_err.RestErr)
}
