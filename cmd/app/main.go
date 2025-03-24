package main

import (
	"context"
	"log"

	"github.com/go-chi/chi/v5"

	"github.com/Ekvo/golang-postgres-chi-api/internal/server"
	s "github.com/Ekvo/golang-postgres-chi-api/internal/source"
	"github.com/Ekvo/golang-postgres-chi-api/internal/transport"
)

func main() {
	ctx := context.Background()
	r := chi.NewRouter()
	connect := server.Init(r, ".env")

	db := s.Init(s.URLParam(".env"))
	defer db.Close()
	base := s.NewDbinstance(db)
	// test connection
	if err := base.CreateTables(ctx); err != nil {
		log.Fatalf("create tables error - %v", err)
	}

	h := transport.NewTransport(r)
	h.Routes(base)

	connect.ListenAndServe(ctx, server.TimeShutServer)
}
