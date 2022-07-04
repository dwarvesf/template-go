package entities

import (
	"context"

	"github.com/dwarvesf/go-template/pkg/model"
	"github.com/dwarvesf/go-template/pkg/monitoring"
	"github.com/dwarvesf/go-template/pkg/util"
)

func (e *entity) GetUserByEmail(ctx context.Context, email string) (*model.User, error) {
	return e.repo.UserRepo().GetUserByEmail(ctx, email)
}

func (e *entity) CreateUser(ctx context.Context, user model.User) (*model.User, error) {
	m := monitoring.FromContext(ctx)
	pwdBytes, err := util.HashPassword(user.Password)
	if err != nil {
		m.Errorf(err, "[entity.CreateUser] HashPassword")
		return nil, err
	}

	user.Password = string(pwdBytes)
	return e.repo.UserRepo().CreateUser(ctx, user)
}
