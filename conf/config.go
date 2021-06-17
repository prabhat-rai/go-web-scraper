package conf

import (
	"go.mongodb.org/mongo-driver/mongo"
	"os"
)

type (
	Config struct {
		AllApps AllApps
		AndroidApps AllAndroidApps
		IosApps AllIosApps
		AllSchedulerConfigs AllSchedulerConfigs
	}
)

func New(client *mongo.Client) *Config {
	// Connect to DB
	dbName := os.Getenv("DB_DATABASE")
	database := client.Database(dbName)

	return &Config{
		AllApps: GetAppsConfig(database, true),
		AndroidApps: GetAndroidAppsViaConfig(),
		IosApps: GetIosAppsViaConfig(),
		AllSchedulerConfigs:GetSchedulerConfigsViaConfig(),
	}
}