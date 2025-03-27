package main

import (
	"context"
	"log"

	"github.com/go-chi/chi/v5"

	"github.com/Ekvo/golang-chi-postgres-api/internal/server"
	"github.com/Ekvo/golang-chi-postgres-api/internal/source"
	"github.com/Ekvo/golang-chi-postgres-api/internal/transport"
)

func main() {
	ctx := context.Background()
	r := chi.NewRouter()
	connect := server.Init(r, ".env")

	db := source.Init(source.URLParam(".env"))
	defer db.Close()
	base := source.NewDbinstance(db)
	// test connection
	if err := base.CreateTables(ctx); err != nil {
		log.Fatalf("create tables error - %v", err)
	}

	h := transport.NewTransport(r)
	h.Routes(base)

	connect.ListenAndServeAndShut(ctx, server.TimeShutServer)
}
