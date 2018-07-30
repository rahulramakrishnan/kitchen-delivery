package config

import (
	"fmt"
	"io/ioutil"
	"log"

	yaml "gopkg.in/yaml.v2"
)

// AppConfig holds the application configuration.
type AppConfig struct {
	ServiceName string     `yaml:"service_name"`
	Databases   Databases  `yaml:"databases"`
	Pickup      Pickup     `yaml:"pickup"`
	WorkerPool  WorkerPool `yaml:"worker_pool"`
	ShelfSpace  ShelfSpace `yaml:"shelf_space"`
}

// LoadConfig loads configuration from yaml files.
func (a *AppConfig) LoadConfig(filePath string) error {
	var configFile string

	if filePath == "" {
		configFile = "config/development.yaml"
	} else {
		configFile = filePath
	}

	// Load configuration file.
	yamlFile, err := ioutil.ReadFile(configFile)
	if err != nil {
		log.Printf("failed to read from yamlFile.Get | err: %s", err)
		return err
	}

	err = yaml.Unmarshal(yamlFile, a)
	if err != nil {
		log.Printf("failed to unmarshal | err: %s", err)
		return err
	}

	return nil
}

// Databases holds database connection information.
type Databases struct {
	MySQL MySQL `yaml:"mysql"`
}

// MySQL holds master and slave SQL connection urls.
type MySQL struct {
	Username string `yaml:"username"`
	Database string `yaml:"database"`
}

// GetConnectionString returns MySQL connection string.
func (m *MySQL) GetConnectionString() string {
	// Fetch required environment variables.
	connectionStr := fmt.Sprintf("%s:@tcp(localhost:3306)/%s?charset=utf8&parseTime=True&loc=Local", m.Username, m.Database)

	return connectionStr
}

// Pickup holds pickup information.
type Pickup struct {
	Mean float64 `yaml:"mean"` // mean for poisson distribution
}

// WorkerPool holds max worker count.
type WorkerPool struct {
	MaxWorkers int `yaml:"max_workers"` // num of max workers.
}

// ShelfSpace holds capacity of each type of shelf.
type ShelfSpace struct {
	Hot      int `yaml:"hot"`
	Cold     int `yaml:"cold"`
	Frozen   int `yaml:"frozen"`
	Overflow int `yaml:"overflow"`
}
