package config

import (
	"fmt"
	"os"
	"regexp"
)

var (
	Settings       = getConfig()
	mongoUriFormat = `^mongodb\:\/\/(?P<username>[_\w]+):(?P<password>[\w]+)@(?P<host>[\.\w]+):(?P<port>\d+)/(?P<database>[_\w]+)$`
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

	config.Port = ":" + os.Getenv("PORT")
	if config.Port == ":" {
		config.Port = ":8080"
	}

	// Set up database connection url
	config.DatabaseUrl = os.Getenv("MONGOLAB_URI")
	if config.DatabaseUrl == "" {
		config.DatabaseUrl = "localhost"
		config.DatabaseName = "wakeup-call-dev"
	} else {
		// TODO: Grab databaseName from mongoURL using regex
		//config.DatabaseName = "wakeup-call-prod"
		regex, _ := regexp.Compile(mongoUriFormat)
		matches := regex.FindAllStringSubmatch(config.DatabaseUrl, -1)[0]
		fmt.Printf("The matches are: %v\n", matches)
		config.DatabaseName = "heroku_app33135020"
	}

	return config
}
