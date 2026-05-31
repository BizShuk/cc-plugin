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
	viper.SetDefault("state.db_path", "~/.config/cc-plugin/state.db")
	viper.SetDefault("retention.max_age_days", 30)
	viper.SetDefault("llm.host", "http://localhost:11434")
	viper.SetDefault("llm.model", "qwen3:14b-q4_K_M")
	viper.SetDefault("sources.claude_mem.db_path", "~/.claude-mem/claude-mem.db")
	viper.SetDefault("sources.gbrain_working.root", "~/brain/working")
	viper.SetDefault("stores.agentmemory.url", "http://localhost:3111/agentmemory/remember")
	viper.SetDefault("stores.mempalace.wing", "main")
	viper.SetDefault("stores.mempalace.temp_dir", "/tmp/mempalace-temp")
}
