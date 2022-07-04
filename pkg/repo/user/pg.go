package user

import (
	"context"

	"github.com/dwarvesf/go-template/pkg/model"
	"gorm.io/gorm"
)

type userRepo struct {
	db *gorm.DB
}

func NewUserRepoPg(db *gorm.DB) Store {
	return &userRepo{
		db: db,
	}
}

func (s userRepo) GetUserByEmail(ctx context.Context, email string) (*model.User, error) {
	var res model.User
	return &res, s.db.Where("lower(email) = lower(?)", email).First(&res).Error
}

func (s userRepo) CreateUser(ctx context.Context, user model.User) (*model.User, error) {
	return &user, s.db.Create(&user).Error
}
