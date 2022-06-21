package repo

import (
	"github.com/dwarvesf/go-template/pkg/repo/user"
)

// TxFunc function to finish a transaction
type TxFunc = func(store Store) error

// Store persistent data interface
type Store interface {
	DoInTransaction(txFunc TxFunc) error
	UserRepo() user.Store
}
