package config

import (
	"github.com/bizshuk/gosdk/config"
	"github.com/spf13/viper"
)

func Init() {
	config.Default(config.WithAppName("cc-plugin"))
	viper.SetDefault("state.db_path", "~/.config/cc-plugin/state.db")
	viper.SetDefault("sources.claude_mem.db_path", "~/.claude-mem/claude-mem.db")
}
