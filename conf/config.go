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

	DbConfig struct {
		Host 		string `mapstructure:"DB_HOST"`
		Port 		string `mapstructure:"DB_PORT"`
		Database 	string `mapstructure:"DB_DATABASE"`
		User  		string `mapstructure:"DB_USER"`
		Password 	string `mapstructure:"DB_PASSWORD"`
	}

	MailConfig struct {
		Host 		string `mapstructure:"MAIL_HOST"`
		Port 		string `mapstructure:"MAIL_PORT"`
		User  		string `mapstructure:"MAIL_USER"`
		Password 	string `mapstructure:"MAIL_PASSWORD"`
		SendMail	string `mapstructure:"SEND_MAIL"`
	}

	ConfigProps struct {
		AppPort      string `mapstructure:"APP_PORT"`
		DbConfig DbConfig `mapstructure:"DB_CONFIG"`
		SchedulerConfigs map[string]string `mapstructure:"SCHEDULER_CONFIGS"`
		MailConfig MailConfig `mapstructure:"MAIL_CONFIG"`
	}

)

func New(client *mongo.Client,configProps ConfigProps) *Config {

	// Connect to DB
	dbName := configProps.DbConfig.Database
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