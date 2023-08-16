package config

import (
	"fmt"
	"os"
	"path/filepath"

	"yellow-jersey/pkg/logs"
	"yellow-jersey/pkg/project"

	"github.com/BurntSushi/toml"
)

type (
	// Config contains config for the application.
	Config struct {
		DB
		Strava
	}

	// DB holds database values in our config.
	DB struct {
		Host     string `toml:"db_host"`
		Port     string `toml:"db_port"`
		Password string `toml:"db_password"`
		Username string `toml:"db_user"`
	}

	// Strava holds config related to accessing the Strava-API.
	Strava struct {
		ClientSecret string `toml:"strava_client_secret"`
		ClientID     int64  `toml:"strava_client_id"`
	}
)

// NewConfig creates a new config struct.
// TODO: Implement secret manager or vault to override secret config vars
func NewConfig() (*Config, error) {
	config := &Config{}
	if _, err := toml.DecodeFile(generatePath(project.Root()), &config); err != nil {
		return nil, err
	}
	logs.Logger.Debug().Msgf("loaded config %+v from %s", config, generatePath(project.Root()))

	return config, nil
}

// ConnectionString returns the postgres url connection.
func (d DB) ConnectionString() string {
	return fmt.Sprintf("mongodb://%s:%s@%s:%s/?maxPoolSize=20&w=majority", d.Username, d.Password, d.Host, d.Port)
}

func environment() string {
	environment := "dev"
	if os.Getenv("ENV") != "" {
		environment = os.Getenv("ENV")
	}

	return environment
}

func generatePath(projectRoot string) string {
	return filepath.Join(projectRoot, "config", fmt.Sprintf("%s.toml", environment()))
}
