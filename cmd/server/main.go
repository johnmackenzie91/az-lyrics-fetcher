package main

import (
	"context"
	"syscall"

	"os"
	"os/signal"

	"time"

	"github.com/johnmackenzie91/azlyrics-fetcher/internal/app"
	"github.com/johnmackenzie91/azlyrics-fetcher/internal/client"
	"github.com/johnmackenzie91/azlyrics-fetcher/internal/http"
)

func main() {

	azClient, _ := client.New()
	a, err := app.New(azClient)

	if err != nil {
		panic(err)
	}

	s := http.NewServer("0.0.0.0:80", a)

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	ctx := context.Background()

	go func() {
		if err := s.Run(ctx); err != nil {
			panic(err)
		}
	}()

	<-done

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer func() {
		// extra handling here
		cancel()
	}()

	if err := s.Shutdown(ctx); err != nil {
	}
}
