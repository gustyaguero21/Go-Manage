package main

import (
	"go-manage/cmd/config"
	"go-manage/internal/router"
	"log"
)

func main() {
	router := router.SetupRouter()

	if err := router.Run(config.Port); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
