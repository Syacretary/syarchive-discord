package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	DiscordToken           string `mapstructure:"DISCORD_TOKEN"`
	OpenRouterAPIKey       string `mapstructure:"OPENROUTER_API_KEY"`
	GoogleSearchAPIKey     string `mapstructure:"GOOGLE_SEARCH_API_KEY"`
	GoogleSearchEngineID   string `mapstructure:"GOOGLE_SEARCH_ENGINE_ID"`
	BotPrefix              string `mapstructure:"BOT_PREFIX"`
	MaxConcurrentDownloads int    `mapstructure:"MAX_CONCURRENT_DOWNLOADS"`
	MaxFileSize            int    `mapstructure:"MAX_FILE_SIZE"`
}

func LoadConfig() (*Config, error) {
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()

	// Default values
	viper.SetDefault("BOT_PREFIX", "/")
	viper.SetDefault("MAX_CONCURRENT_DOWNLOADS", 3)
	viper.SetDefault("MAX_FILE_SIZE", 100)

	if err := viper.ReadInConfig(); err != nil {
		// Jika file .env tidak ditemukan, kita tetap bisa menggunakan environment variables
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, err
		}
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}

	return &config, nil
}