package config

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	yaml "gopkg.in/yaml.v2"
)

// AppConfig holds the application configuration.
type AppConfig struct {
	ServiceName string    `yaml:"service_name"`
	Databases   Databases `yaml:"databases"`
}

// LoadConfig loads configuration from yaml files.
func (a *AppConfig) LoadConfig() error {
	configFile := "config/development.yaml"

	// Load configuration file.
	yamlFile, err := ioutil.ReadFile(configFile)
	if err != nil {
		log.Printf("failed to read from yamlFile.Get | err: %+v", err)
		return err
	}

	err = yaml.Unmarshal(yamlFile, a)
	if err != nil {
		log.Printf("failed to unmarshal | err: %v", err)
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
	Password string `yaml:"password"`
	Database string `yaml:"database"`
}

// GetConnectionString returns MySQL connection string.
func (m *MySQL) GetConnectionString() string {
	// Fetch required environment variables.
	password := os.Getenv("MYSQL_PASSWORD")
	connectionStr := fmt.Sprintf("%s:@tcp(localhost:3306)/%s?charset=utf8&parseTime=True&loc=Local", m.Username, m.Database)

	return connectionStr
}
