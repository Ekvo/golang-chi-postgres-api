// source - describes a container for pointe sql.DB
package source

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type Dbinstance struct {
	db *sql.DB
}

func NewDbinstance(db *sql.DB) *Dbinstance {
	return &Dbinstance{db: db}
}

func Init(dataSourceName string) *sql.DB {
	db, err := sql.Open("postgres", dataSourceName)
	if err != nil {
		log.Fatalf("sql.Open - %v", err)
	}
	return db
}

// URLParam - get from .env file and create dataSourceName for 'sql.Opne(,dataSourceName)'
//
//	postgres://jack:secret@pg.example.com:5432/mydb?sslmode=verify-ca&pool_max_conns=10&pool_max_conn_lifetime=1h30m
func URLParam(path string) string {
	if err := godotenv.Load(path); err != nil {
		log.Fatalf("source - no .env data: %v", err)
	}
	dsn := fmt.Sprintf(
		`postgres://%s:%s@%s:%s/%s?sslmode=%s`,
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_SSLMODE"),
	)
	return dsn
}
