package user_entity

import (
	"context"

	"github.com/valdir-alves3000/postgraduate-challenges-go-expert/labs-auction-goexpert/internal/internal_error"
)

type User struct {
	Id   string
	Name string
}

type UserRepositoryInterface interface {
	FindUserById(
		ctx context.Context, userId string) (*User, *internal_error.InternalError)
}
