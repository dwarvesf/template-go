package repo

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/dwarvesf/go-template/pkg/config"
	user "github.com/dwarvesf/go-template/pkg/repo/user"
	_ "github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// NewPostgresStore postgres init by gorm
func NewPostgresStore(cfg *config.Config) (Store, func() error) {
	ds := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.DBUser, cfg.DBPass,
		cfg.DBHost, cfg.DBPort, cfg.DBName,
	)
	conn, err := sql.Open("postgres", ds)
	if err != nil {
		panic(err)
	}

	db, err := gorm.Open(postgres.New(
		postgres.Config{Conn: conn}),
		&gorm.Config{
			// TODO: replace with monitoring pkg
			Logger: logger.New(
				log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
				logger.Config{
					SlowThreshold:             time.Second,   // Slow SQL threshold
					LogLevel:                  logger.Silent, // Log level
					IgnoreRecordNotFoundError: true,          // Ignore ErrRecordNotFound error for logger
					Colorful:                  true,          // Disable color
				},
			),
		})
	if err != nil {
		panic(err)
	}

	return newRepo(db), conn.Close
}

type repoIpl struct {
	db   *gorm.DB
	user user.Store
}

func newRepo(db *gorm.DB) Store {
	return &repoIpl{
		user: user.NewUserRepoPg(db),
	}
}

func (r repoIpl) UserRepo() user.Store {
	return r.user
}

func (r repoIpl) DoInTransaction(txFunc TxFunc) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		s := newRepo(tx)
		return txFunc(s)
	})
}
