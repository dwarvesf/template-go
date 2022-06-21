package model

import (
	"github.com/google/uuid"
	"github.com/volatiletech/null"
)

type User struct {
	ID          uuid.NullUUID `json:"id" gorm:"default:uuid_generate_v4()"`
	Email       string        `json:"email"`
	Password    string        `json:"-"`
	CreatedAt   null.Time     `json:"created_at" gorm:"default:now()"`
	UpdatedAt   null.Time     `json:"updated_at" gorm:"default:now()"`
	AccessToken string        `json:"access_token" gorm:"-"`
}
