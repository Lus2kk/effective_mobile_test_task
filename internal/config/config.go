package config

import (
	"errors"
	"log/slog"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	AppEnv   string         `env:"APP_ENV"`
	Server   ServerConfig	
	Database DatabaseConfig  
}

type DatabaseConfig struct {
	User     string `env:"POSTGRES_USER"`
	Password string `env:"POSTGRES_PASSWORD"`
	Addr     string `env:"POSTGRES_ADDR"`
	DB       string `env:"POSTGRES_DB"`
}
type ServerConfig struct {
	Port string `env:"SERVER_PORT"`
	Host string `env:"SERVER_HOST"`
}

func MustLoad() (*Config, error) {
    var cfg Config
    
    if _, err := os.Stat(".env"); err == nil {
        err = cleanenv.ReadConfig(".env", &cfg)
        if err != nil {
            slog.Error("Failed to read config", "error", err)
            return nil, errors.New("failed to load config")
        }
    } else {
        err = cleanenv.ReadEnv(&cfg)
        if err != nil {
            slog.Error("Failed to read environment variables", "error", err)
            return nil, errors.New("failed to load config")
        }
    }
    
    return &cfg, nil
}
