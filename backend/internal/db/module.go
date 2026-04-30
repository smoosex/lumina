package db

import (
	"database/sql"

	"github.com/smoosex/lumina/backend/internal/config"
	"go.uber.org/fx"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var Module = fx.Module("db", fx.Provide(Open))

func Open(lc fx.Lifecycle, cfg config.Config) (*gorm.DB, error) {
	conn, err := gorm.Open(postgres.Open(cfg.Database.DSN), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	sqlDB, err := conn.DB()
	if err != nil {
		return nil, err
	}

	sqlDB.SetMaxOpenConns(10)
	sqlDB.SetMaxIdleConns(5)

	lc.Append(fx.StopHook(func() error {
		return closeSQL(sqlDB)
	}))

	return conn, nil
}

func closeSQL(db *sql.DB) error {
	if db == nil {
		return nil
	}
	return db.Close()
}
