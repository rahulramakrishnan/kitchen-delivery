package main

import (
	"log"
	"net/http"

	"github.com/gomodule/redigo/redis"
	"github.com/kitchen-delivery/config"
	"github.com/kitchen-delivery/entity"
	"github.com/kitchen-delivery/handler"
	"github.com/kitchen-delivery/job"
	"github.com/kitchen-delivery/service"
	"github.com/kitchen-delivery/service/repository"

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

	////////////////////////////////////////
	// Storage Initialization
	////////////////////////////////////////

	// Open connection to MySQL instance.
	db, err := gorm.Open("mysql", cfg.Databases.MySQL.GetConnectionString())
	if err != nil {
		log.Fatalf("Failed to connect to mysql database %+v", err)
	}
	defer db.Close()

	// Open connection to Redis instance.
	// Use this as a first in first out queue.
	// Instead of the channel.
	redisConn, err := redis.Dial("tcp", ":6379")
	if err != nil {
		log.Fatalf("Failed to connect to redis instance %+v", err)
	}
	defer redisConn.Close()

	////////////////////////////////////////
	// Service Initialization
	////////////////////////////////////////
	repositories := repository.InitializeRepositories(db)
	services := service.InitializeServices(cfg, repositories)

	////////////////////////////////////////
	// Local Queue Initialization
	////////////////////////////////////////
	orderQueue := make(chan *entity.Order)      // TODO: Replace w/ Redis
	shelfOrderQueue := make(chan *entity.Order) // TODO: Replace w/ Redis
	queues := entity.Queues{
		OrderQueue: orderQueue,
		ShelfQueue: shelfOrderQueue,
	}

	////////////////////////////////////////
	// Job & Worker Initialization
	////////////////////////////////////////
	jobs := job.InitializeJobs(cfg, services, queues)

	// Spawn workers to pull orders off of order queue
	// as orders come in.
	go jobs.Order.HandleIncomingOrders()

	// Spawn thread to remove expired orders.
	go jobs.Order.RemoveExpiredOrders()

	////////////////////////////////////////
	// Handler Initialization
	////////////////////////////////////////
	handlers, err := handler.NewHandlers(cfg, services, orderQueue)
	if err != nil {
		log.Fatalf("Failed to initialize handlers - err: %+v", err)
	}

	////////////////////////////////////////
	// HTTP Route Initialization
	////////////////////////////////////////

	// Register service health and simulation routes.
	http.HandleFunc("/health", handlers.Health.CheckHealth)
	http.HandleFunc("/health/simulate", handlers.Health.Simulate)

	// Register order routes.
	http.HandleFunc("/order", handlers.Order.HandleOrder)

	log.Print("Kitchen Delivery online ....")

	// Mount server and listen on HTTP port.
	http.ListenAndServe(":8080", nil)

	// Block indefinitely to keep server alive.
	switch {

	}
}
