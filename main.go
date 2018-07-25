package main

import (
	"log"
	"net/http"

	"github.com/kitchen-delivery/config"
	"github.com/kitchen-delivery/handler"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

func main() {
	log.Print("Starting Kitchen Delivery ....")

	// Load application configuration.
	cfg := config.AppConfig{}
	err := cfg.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration from yaml files - err: %+v", err)
	}

	// Initialize a MySQL DB connection.
	connectionString := cfg.Databases.MySQL.GetConnectionString()

	db, err := gorm.Open("mysql", connectionString)
	if err != nil {
		log.Fatalf("Failed to connect to mysql database - err: %+v", err)
	}
	defer db.Close()

	////////////////////////////////////////
	// Handler Initialization
	////////////////////////////////////////
	handlers, err := handler.NewHandlers(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize handlers - err: %+v", err)
	}

	////////////////////////////////////////
	// HTTP Route Initialization
	////////////////////////////////////////

	// Register service health routes.
	http.HandleFunc("/health", handlers.Health.CheckHealth)

	// Register order routes.
	http.HandleFunc("/order", handlers.Order.HandleOrder)

	log.Print("Kitchen Delivery online ....")

	// Mount server and listen on HTTP port.
	http.ListenAndServe(":8080", nil)

	// Block indefinitely to keep server alive.
	switch {

	}
}
