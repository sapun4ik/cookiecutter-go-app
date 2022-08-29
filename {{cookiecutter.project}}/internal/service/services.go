package service

import (
	"{{cookiecutter.module_path}}/internal/repository"
)

type Services struct {
}

type Deps struct {
	Repos *repository.Repositories
}

func NewServices(deps Deps) *Services {
	return &Services{}
}
