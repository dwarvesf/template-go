package testutil

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/dwarvesf/go-template/pkg/config"
	pgx "github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/stdlib"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// appDB caches a pg connection for reuse
var appDB *sql.DB

// WithTxDB provides callback with a `pg.BeginnerExecutor` for running pg related tests
// where the `pg.BeginnerExecutor` is actually powered by a pg transaction
// and will be rolled back (so no data is actually written into pg)
func WithTxDB(t *testing.T, cfg config.Config, callback func(*gorm.DB)) {
	if appDB == nil {
		var err error

		dbConnCfg, err := pgx.ParseConfig(cfg.GetDBURL())
		require.NoError(t, err)
		dbConnCfg.LogLevel = pgx.LogLevelDebug
		connStr := stdlib.RegisterConnConfig(dbConnCfg)
		appDB, err = sql.Open("pgx", connStr)
		appDB.SetMaxOpenConns(50)
		appDB.SetConnMaxLifetime(30 * time.Minute)
		require.NoError(t, err)
	}

	tx, err := appDB.BeginTx(context.Background(), nil)
	require.NoError(t, err)

	// explicitly hardcoded to rollback via `defer` because
	// in case `callback` does a t.Fatal we must still run tx.Rollback
	defer tx.Rollback()

	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: tx,
	}), &gorm.Config{})
	require.NoError(t, err)

	callback(gormDB)
}
