package conf

type (
	Config struct {
		AndroidApps AllAndroidApps
		IosApps AllIosApps
	}
)

func New() *Config {
	return &Config{
		AndroidApps: GetAndroidAppsConfig(),
		IosApps: GetIosAppsConfig(),
	}
}