package main

import (
	"log"

	"github.com/kelseyhightower/envconfig"
)

func main() {

	var config Config
	if err := envconfig.Process("", &config); err != nil {
		log.Fatalf("Failed to load config: %w\n", err)
	}

	db, err := InitDB(&config)
	if err != nil {
		log.Fatalf("Failed to load database: %w\n", err)
	}

	router := GetRouter(db, &config)

	if err := router.Run(":" + config.Port); err != nil {
		log.Fatalf("Failed to run application: %w\n", err)
	}
}
