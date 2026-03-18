package main

import (
	"log"

	"github.com/pavandhadge/ssh_portfolio_opensource/internal/server"
)

func main() {
	if err := server.Run(); err != nil {
		log.Printf("server shutdown error: %v", err)
	}
}
