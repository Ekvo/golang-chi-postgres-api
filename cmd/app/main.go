package main

import (
	"context"
	"log"

	"github.com/go-chi/chi/v5"

	"github.com/Ekvo/golang-chi-postgres-api/internal/config"
	"github.com/Ekvo/golang-chi-postgres-api/internal/server"
	"github.com/Ekvo/golang-chi-postgres-api/internal/source"
	"github.com/Ekvo/golang-chi-postgres-api/internal/transport"
)

func main() {
	cfg, err := config.NewConfig(".env", false)
	if err != nil {
		log.Fatalf("main: error - %v", err)
	}

	db, err := source.Init(cfg)
	if err != nil {
		log.Fatalf("main: db error - %v", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			log.Printf("main: db.Close error - %v", err)
		}
	}()
	ctx := context.Background()
	base := source.NewDbinstance(db)
	if err := base.CreateTables(ctx); err != nil {
		log.Fatalf("create tables error - %v", err)
	}
	r := chi.NewRouter()
	connect := server.Init(cfg, r)
	transport.NewTransport(r).Routes(base)

	if err := connect.ListenAndServeAndShut(ctx, server.TimeShutServer); err != nil {
		log.Fatalf("main: server error - %v", err)
	}
}
