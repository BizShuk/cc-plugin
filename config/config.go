package config

import (
	_ "embed"

	"github.com/bizshuk/gosdk/config"
)

//go:embed default_settings.json
var defaultSettingJSON string

func Init() {
	config.Default(
		config.WithAppName("cc-plugin"),
		config.WithDefaultValue(defaultSettingJSON),
	)
}
