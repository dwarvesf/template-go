package entities

import (
	"context"

	"github.com/dwarvesf/go-template/pkg/config"
	"github.com/dwarvesf/go-template/pkg/model"
	"github.com/dwarvesf/go-template/pkg/repo"
)

type entity struct {
	cfg  config.Config
	repo repo.Store
	// log  logger.Log
}
type Entity interface {
	LoginUser(ctx context.Context, email, pwd string) (*model.User, error)
	GetUserByEmail(ctx context.Context, email string) (*model.User, error)
	CreateUser(ctx context.Context, user model.User) (*model.User, error)
}

// l logger.Log
func New(cfg config.Config, r repo.Store) Entity {
	return &entity{
		cfg:  cfg,
		repo: r,
	}
}
