package main

import (
	"ProjectGolang/internal/config"
	"ProjectGolang/pkg/log"
	"github.com/joho/godotenv"
)

func main() {
	logger := log.NewLogger()

	if err := godotenv.Load(); err != nil {
		logger.Fatal("Error loading .env file")
	}

	fiber := config.NewFiber(logger)
	validator := config.NewValidator()
	rest, err := config.NewServer(fiber, logger, validator)
	if err != nil {
		logger.Fatal(err)
	}

	rest.RegisterHandler()

	if err := rest.Run(); err != nil {
		logger.Fatal(err)
	}
}
