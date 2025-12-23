package config

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	MySQL MySQLConfig
	Log   LogConfig
	Redis RedisConfig
}
type MySQLConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
}
type LogConfig struct {
	Level      string
	Filename   string
	MaxSize    int `mapstructure:"max_size"`
	MaxAge     int `mapstructure:"max_age"`
	MaxBackups int `mapstructure:"max_backups"`
	Compress   bool
}
type RedisConfig struct {
	Addr     string
	Password string
	DB       int
	PoolSize int `mapstructure:"pool_size"`
}

var GlobalConfig *Config

func InitConfig() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error read config file, %s", err)
	}
	if err := viper.Unmarshal(&GlobalConfig); err != nil {
		log.Fatalf("Ubable todecode into struct, %v", err)
	}
	log.Println("配置加载成功")
}
