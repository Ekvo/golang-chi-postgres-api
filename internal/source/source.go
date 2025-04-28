// source - describes a container for pointe sql.DB
package source

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"

	"github.com/Ekvo/golang-chi-postgres-api/internal/config"
)

type Dbinstance struct {
	db *sql.DB
}

func NewDbinstance(db *sql.DB) *Dbinstance {
	return &Dbinstance{db: db}
}

func Init(cfg *config.Config) (*sql.DB, error) {
	db, err := sql.Open("postgres", dbURL(cfg))
	if err != nil {
		return nil, fmt.Errorf("source: sql.Open error - %w", err)
	}
	if err := db.Ping(); err != nil {
		go func() {
			if err := db.Close(); err != nil {
				log.Printf("source: DB.Close error - %v", err)
			}
		}()
		return nil, fmt.Errorf("source: DB.Ping error - %w", err)
	}
	return db, nil
}

// URLParam - get from .env file and create dataSourceName for 'sql.Opne(,dataSourceName)'
//
//	postgres://jack:secret@pg.example.com:5432/mydb?sslmode=verify-ca&pool_max_conns=10&pool_max_conn_lifetime=1h30m
func dbURL(cfg *config.Config) string {
	dsn := fmt.Sprintf(
		`postgres://%s:%s@%s:%s/%s?sslmode=%s`,
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBHost,
		cfg.DBPort,
		cfg.DBName,
		cfg.DBSSLMode,
	)
	return dsn
}
