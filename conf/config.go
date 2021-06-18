package conf

import (
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/mongo"
)


type (
	Config struct {
		AllApps     AllApps
		ConfigProps ConfigProps

	}
	ConfigProps struct {
		APP_PORT      string `mapstructure:"APP_PORT"`
		DB_HOST     string `mapstructure:"DB_HOST"`
		DB_PORT     string `mapstructure:"DB_PORT"`
		DB_DATABASE string `mapstructure:"DB_DATABASE"`
		DB_USER  string `mapstructure:"DB_USER"`
		DB_PASSWORD string `mapstructure:"DB_PASSWORD"`
		SCHEDULER_CONFIGS map[string]string `mapstructure:"SCHEDULER_CONFIGS"`
	}

)

func New(client *mongo.Client,configProps ConfigProps) *Config {

	// Connect to DB
	dbName := configProps.DB_DATABASE
	database := client.Database(dbName)
	return &Config{
		AllApps: GetAppsConfig(database, true),
		ConfigProps: configProps,
	}
}

// LoadConfig reads configuration from file or environment variables.
func LoadConfig(path string) (config ConfigProps, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("config")
	viper.SetConfigType("yml")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}