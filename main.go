package main

import (
	"log"

	"github.com/kitchen-delivery/config"

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

	log.Print("Kitchen Delivery online ....")

	// Block indefinitely to keep server alive.
	switch {

	}
}
