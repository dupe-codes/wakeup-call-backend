package config

import (
    "os"
    //"regexp"
)

var (
	Settings = getConfig()
	//mongoUriFormat = "^mongodb\:\/\/(?P<username>[_\w]+):(?P<password>[\w]+)@(?P<host>[\.\w]+):(?P<port>\d+)/(?P<database>[_\w]+)$"
)

type Config struct {
	ProjectName  string
	Port         string
	DatabaseUrl  string
	DatabaseName string
	DbUsername   string
	DbPassword   string
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
	    //config.DatabaseName = "wakeup-call-prod"
	    config.DatabaseName = "heroku_app33135020"
	}
	
	// TODO: Generate this programatically
	// Credentials for heroku mongodb server
	// regex, _ := regexp.Compile(mongoUriFormat)

	return config
}
