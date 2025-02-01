package main

import (
	"go-manage/cmd/config"
	"go-manage/internal/router"
)

func main() {
	router := router.SetupRouter()

	router.Run(config.Port)
}
