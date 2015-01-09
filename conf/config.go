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
	config.Port = ":8080"


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
