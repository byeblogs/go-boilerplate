package database

import (
	"fmt"
	"log"

	"github.com/byeblogs/go-boilerplate/pkg/config"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/jmoiron/sqlx"
)

type DB struct{ *sqlx.DB }

var defaultDB = &DB{}

func (db *DB) connect(cfg *config.DB) (err error) {
	// Build DSN with safe defaults already applied by config.LoadDBCfg
	dsn := config.BuildPostgresDSN(cfg)

	// Helpful one-line (redacted) to confirm envs arrived at runtime
	log.Printf("DB connect → host=%s port=%d user=%s db=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.Name, cfg.SslMode)

	db.DB, err = sqlx.Connect("pgx", dsn)
	if err != nil {
		return fmt.Errorf("connect failed: %w", err)
	}

	db.SetMaxOpenConns(cfg.MaxOpenConn)
	db.SetMaxIdleConns(cfg.MaxIdleConn)
	db.SetConnMaxLifetime(cfg.MaxConnLifetime)

	if err := db.Ping(); err != nil {
		defer db.Close()
		return fmt.Errorf("can't ping database: %w", err)
	}
	return nil
}

func GetDB() *DB       { return defaultDB }
func ConnectDB() error { return defaultDB.connect(config.DBCfg()) }
