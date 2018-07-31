package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/kitchen-delivery/config"
	"github.com/kitchen-delivery/entity"
	"github.com/kitchen-delivery/handler"
	"github.com/kitchen-delivery/job"
	"github.com/kitchen-delivery/service"
	"github.com/kitchen-delivery/service/repository"

	"github.com/gomodule/redigo/redis"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

func main() {
	log.Print("Starting Kitchen Delivery ....")

	// Load application configuration.
	cfg := config.AppConfig{}
	err := cfg.LoadConfig("config/development.yaml")
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
	// TODO: Move this to a configuration.
	redisPool := &redis.Pool{
		MaxIdle:     cfg.Redis.MaxIdle,
		MaxActive:   cfg.Redis.MaxActive,
		IdleTimeout: time.Duration(cfg.Redis.IdleTimeout) * time.Second,
		Wait:        cfg.Redis.Wait,
		Dial: func() (redis.Conn, error) {
			redisHost := fmt.Sprintf(":%d", cfg.Redis.Port)
			redisConn, err := redis.Dial("tcp", redisHost)
			if err != nil {
				return nil, err
			}

			return redisConn, nil
		},
	}

	////////////////////////////////////////
	// Service Initialization
	////////////////////////////////////////
	repositories := repository.InitializeRepositories(db)
	services := service.InitializeServices(cfg, repositories)

	////////////////////////////////////////
	// Local Queue Initialization
	////////////////////////////////////////
	queues := &entity.Queues{
		// We pass redis pool by reference
		// as it contains mutex lock.
		Order: entity.Queue{
			Name: "Order",
			Pool: redisPool,
		},
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
	handlers, err := handler.NewHandlers(cfg, services, queues)
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
