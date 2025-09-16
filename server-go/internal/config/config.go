package config

import (
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Server   ServerConfig   `mapstructure:"server"`
	Database DatabaseConfig `mapstructure:"database"`
	JWT      JWTConfig      `mapstructure:"jwt"`
	Upload   UploadConfig   `mapstructure:"upload"`
	Log      LogConfig      `mapstructure:"log"`
}

type ServerConfig struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
}

type DatabaseConfig struct {
	URI    string `mapstructure:"uri"`
	DBName string `mapstructure:"db_name"`
}

type JWTConfig struct {
	Secret         string        `mapstructure:"secret"`
	TokenExpires   time.Duration `mapstructure:"token_expires"`
	RefreshExpires time.Duration `mapstructure:"refresh_expires"`
}

type UploadConfig struct {
	Path    string `mapstructure:"path"`
	BaseURL string `mapstructure:"base_url"`
}

type LogConfig struct {
	Level  string `mapstructure:"level"`
	Format string `mapstructure:"format"`
}

func Load() (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./configs")
	viper.AddConfigPath(".")

	// 设置默认值
	setDefaults()

	// 环境变量覆盖
	viper.AutomaticEnv()

	// 读取配置文件
	if err := viper.ReadInConfig(); err != nil {
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

func setDefaults() {
	// 服务器配置
	viper.SetDefault("server.host", "0.0.0.0")
	viper.SetDefault("server.port", 8080)

	// 数据库配置
	viper.SetDefault("database.uri", "mongodb://localhost:27017")
	viper.SetDefault("database.db_name", "blog")

	// JWT配置
	viper.SetDefault("jwt.secret", "your-secret-key")
	viper.SetDefault("jwt.token_expires", 24*time.Hour)
	viper.SetDefault("jwt.refresh_expires", 7*24*time.Hour)

	// 上传配置
	viper.SetDefault("upload.path", "./public")
	viper.SetDefault("upload.base_url", "http://localhost:8080")

	// 日志配置
	viper.SetDefault("log.level", "info")
	viper.SetDefault("log.format", "json")
}