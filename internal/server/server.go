// server  - describes a container for pointer http.Server
package server

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
)

const TimeShutServer = 10 * time.Second

// Connect - wrapper http.Server
type Connect struct {
	*http.Server
}

func NewServer(srv *http.Server) *Connect {
	return &Connect{srv}
}

func Init(r http.Handler, envPath string) *Connect {
	err := godotenv.Load(envPath)
	if err != nil {
		log.Fatalf("server - no .env data: %v", err)
	}
	addr := os.Getenv("SRV_ADDR")
	//timeRead, err := strconv.Atoi(os.Getenv("SRV_TIME_READ"))
	//if err != nil {
	//	log.Fatalf("server - .env data incorrect: %v", err)
	//}
	//timeWrite, err := strconv.Atoi(os.Getenv("SRV_TIME_WRITE"))
	//if err != nil {
	//	log.Fatalf("server - .env data incorrect: %v", err)
	//}
	srv := &http.Server{
		Addr:    addr,
		Handler: r,
		//ReadTimeout:  time.Duration(timeRead),
		//WriteTimeout: time.Duration(timeWrite),
	}
	return NewServer(srv)
}

func (c *Connect) ListenAndServeAndShut(ctx context.Context, timeShut time.Duration) {

	go func() {
		if err := c.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("HTTP server error: %v", err)
		}
		log.Println("Stopped serving new connections.")
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	shutdownCtx, shutdownRelease := context.WithTimeout(ctx, timeShut)
	defer shutdownRelease()

	if err := c.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("HTTP shutdown error: %v", err)
	}
	log.Println("Graceful shutdown complete.")
}
