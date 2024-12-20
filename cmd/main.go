package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"tts-poc-service/server"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	srv := server.NewServer(ctx)
	ctx = srv.HandleShutdown(ctx)
	go func() {
		// service connections
		if err := srv.Start(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("error start server: %s\n", err)
		}
	}()
	<-ctx.Done()
}
