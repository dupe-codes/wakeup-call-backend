package config

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

	config.DatabaseUrl = "localhost"
	config.DatabaseName = "wakeup-call-dev"

	return config
}
