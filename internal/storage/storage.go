package storage

import (
	"github.com/lakshsetia/learn-RESTAPI/internal/types"
)

type Storage interface {
	GetUsers() ([]types.User, error)
	CreateUser(name string, email string, age int) error
	GetUserById(id int) (types.User, error)
	UpdateUserById(id int, name string, email string, age int) error
	DeleteUserById(id int) error
}