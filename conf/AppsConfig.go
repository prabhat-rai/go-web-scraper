package conf

import (
	"encoding/json"
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
		AppId string 	`json:"app_id"`
	}

	AllAndroidApps struct {
		Apps []AndroidAppConfig `json:"apps"`
	}

	AllIosApps struct {
		Apps []IosAppConfig `json:"apps"`
	}
)

func GetAndroidAppsConfig() AllAndroidApps {
	// Currently loading from ENV, we can move it to DB
	var androidApps AllAndroidApps
	androidAppConfigs := os.Getenv("ANDROID_APPS")
	err := json.Unmarshal([]byte(androidAppConfigs), &androidApps)

	if err != nil {
		log.Fatal("COULD NOT INTERPRET CONFIG")
	}

	return androidApps
}

func GetIosAppsConfig() AllIosApps {
	// Currently loading from ENV, we can move it to DB
	var iosApps AllIosApps
	iosAppConfigs := os.Getenv("IOS_APPS")
	err := json.Unmarshal([]byte(iosAppConfigs), &iosApps)

	if err != nil {
		log.Fatal("COULD NOT INTERPRET CONFIG")
	}

	return iosApps
}
