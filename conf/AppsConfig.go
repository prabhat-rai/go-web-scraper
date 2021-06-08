package conf

import (
	"context"
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

	AppConfigs struct {
		Name 		string 	`json:"name" bson:"name"`
		GoogleAppId string 	`json:"google_app_id" bson:"google_app_id"`
		IosAppId 	string 	`json:"ios_app_id" bson:"ios_app_id"`
	}

	AllApps struct {
		Apps []AppConfigs `json:"all_apps"`
	}
)

func GetAppsConfig(db *mongo.Database) AllApps {

	var allApps AllApps
	appCollection := db.Collection("apps")
	dbContext := context.TODO()
	app := AppConfigs{}

	// Get All matching Records
	filter := bson.D{{"active", true}}
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
