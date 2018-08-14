package add

import (
	"context"

	"github.com/jinzhu/gorm"
)

type pgStore struct {
	db *gorm.DB
}

// NewPGStore create new project store
func NewPGStore(db *gorm.DB) Service {
	return &pgStore{
		db: db,
	}
}

// Add just do a plus with 2 vars (X+Y) for the sake of demonstration, in reallity
// you would want to execute a DB query/transaction here
func (s *pgStore) Add(ctx context.Context, arg *Add) (int, error) {
	return arg.X + arg.Y, nil
}
