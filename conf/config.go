package config

import (
    "os"
)

var (
	Settings = getConfig()
)

type Config struct {
	ProjectName  string
	Port         string
	DatabaseUrl  string
	DatabaseName string
}

// getConfig sets all relevant project settings
func getConfig() *Config {
	config := new(Config)

	config.ProjectName = "Wakeup_Call"

    port := os.Getenv("PORT")
    if port == "" {
        config.Port = ":8080"
    } else {
        config.Port = ":" + port
    }

    // Set up database connection url
    databaseUrl := os.Getenv("MONGOLAB_URI")
    if databaseUrl == "" {
	    config.DatabaseUrl = "localhost"
        config.DatabaseName = "wakeup-call-dev"
	} else {
	    config.DatabaseUrl = databaseUrl
	    config.DatabaseName = "wakeup-call-prod"
	}

	return config
}
