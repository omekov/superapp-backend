package main

import (
	"log"

	"github.com/omekov/superapp-backend/internal/salecar/server"
)

func main() {
	if err := server.Run(); err != nil {
		log.Fatal(err)
	}
}
