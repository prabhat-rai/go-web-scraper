package conf

import (
	"context"
	"echoApp/model"
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"os"
)

type (
	AndroidAppConfig struct {
		Name string 		`json:"name"`
		GoogleAppId string 	`json:"google_app_id"`
	}
	IosAppConfig struct {
		Name string 	`json:"name"`
		IosAppId string 	`json:"ios_app_id"`
	}

	AllAndroidApps struct {
		Apps []AndroidAppConfig `json:"apps"`
	}

	AllIosApps struct {
		Apps []IosAppConfig `json:"apps"`
	}

	// This is group of model Apps
	AllApps struct {
		Apps []model.Apps `json:"all_apps"`
	}

	SchedulerConfigs struct {
		Concept string `json:"name"`
		Cronexpression string `json:"cronexpression"`
	}
	AllSchedulerConfigs struct {
		SchedulerConfigs []SchedulerConfigs `json:"schedulers"`
	}
)

func GetAppsConfig(db *mongo.Database, onlyActiveRecords bool) AllApps {

	var allApps AllApps
	appCollection := db.Collection("apps")
	dbContext := context.TODO()
	app := model.Apps{}

	filter := bson.D{{}}


	// Get All matching Records
	if onlyActiveRecords {
		filter = bson.D{{"active", true}}
	}

	findOpts := options.Find().SetProjection(bson.M{"ID": 0})
	cursor, err := appCollection.Find(dbContext, filter, findOpts)
	if err != nil {
		return allApps
	}
	for cursor.Next(context.TODO()) {
		err := cursor.Decode(&app)

		if err != nil {
			log.Fatal(err)
		}

		allApps.Apps = append(allApps.Apps, app)
	}

	return allApps
}

func GetAndroidAppsViaConfig() AllAndroidApps {
	var androidApps AllAndroidApps
	androidAppConfigs := os.Getenv("ANDROID_APPS")
	err := json.Unmarshal([]byte(androidAppConfigs), &androidApps)

	if err != nil {
		log.Fatal("COULD NOT INTERPRET CONFIG")
	}

	return androidApps
}


func GetIosAppsViaConfig() AllIosApps {
	var iosApps AllIosApps
	iosAppConfigs := os.Getenv("IOS_APPS")
	err := json.Unmarshal([]byte(iosAppConfigs), &iosApps)

	if err != nil {
		log.Fatal("COULD NOT INTERPRET CONFIG")
	}

	return iosApps
}

func GetSchedulerConfigsViaConfig() AllSchedulerConfigs {
	var schedulerConfigs AllSchedulerConfigs
	schedulerConfigsProp := os.Getenv("SCHEDULER_CONFIGS")
	err := json.Unmarshal([]byte(schedulerConfigsProp), &schedulerConfigs)

	if err != nil {
		log.Fatal("COULD NOT INTERPRET CONFIG")
	}

	return schedulerConfigs
}