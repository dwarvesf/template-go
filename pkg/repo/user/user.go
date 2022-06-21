package user

import (
	"context"

	"github.com/dwarvesf/go-template/pkg/model"
)

type Store interface {
	GetUserByEmail(ctx context.Context, email string) (*model.User, error)
	CreateUser(ctx context.Context, user model.User) (*model.User, error)
}
