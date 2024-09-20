package main

import (
	"context"
	"log"
	"os"
	"os/signal"

	"github.com/omekov/dubaicarkzv2/internal/app"
)

func main() {
	if err := runApp(); err != nil {
		log.Fatal(err)
	}
}

func runApp() error {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	return app.Run(ctx)
}
