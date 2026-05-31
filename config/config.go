package config

import (
	_ "embed"

	"github.com/bizshuk/gosdk/config"
	"github.com/spf13/viper"
)

//go:embed default_settings.json
var defaultSettingJSON string

func Init() {
	config.Default(
		config.WithAppName("cc-plugin"),
		config.WithDefaultValue(defaultSettingJSON),
	)
	viper.SetDefault("retention.max_age_days", 30)
	viper.SetDefault("state.db_path", "~/.config/cc-plugin/state.db")
	viper.SetDefault("stores.agentmemory.url", "http://localhost:3111/agentmemory/remember")
	viper.SetDefault("sources.gbrain_working.root", "~/brain/working")
}
