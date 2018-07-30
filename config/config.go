package config

import (
	"fmt"
	"io/ioutil"
	"log"

	yaml "gopkg.in/yaml.v2"
)

// AppConfig holds the application configuration.
type AppConfig struct {
	ServiceName string    `yaml:"service_name"`
	Databases   Databases `yaml:"databases"`
	Pickup      Pickup    `yaml:"pickup"`
}

// LoadConfig loads configuration from yaml files.
func (a *AppConfig) LoadConfig() error {
	configFile := "config/development.yaml"

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
