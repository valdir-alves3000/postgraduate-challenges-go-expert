package database

import "github.com/valdir-alves3000/postgraduate-challenges-go-expert/APIs/internal/entity"

type UserInterface interface {
	Create(user *entity.User) error
	FindByEmail(email string) (*entity.User, error)
}
