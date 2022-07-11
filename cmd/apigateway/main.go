package main

import (
	"flag"
	"log"

	"github.com/omekov/superapp-backend/internal/apigateway/server"
)

var (
	configPath string
	port       string
)

func main() {
	flagConfigPath()
	flag.Parse()
	if err := server.Run(port, configPath); err != nil {
		log.Fatal(err)
	}
}

func flagConfigPath() {
	flag.StringVar(&configPath, "config-path", "configs/config.yaml", "Path to config files")
	flag.StringVar(&port, "port", ":4041", "tcp port")
}
