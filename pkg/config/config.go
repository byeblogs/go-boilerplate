package config

import (
	"log"
	"os"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

// LoadAllConfigs sets various configs.
// In prod/docker, we read from real environment vars.
// In local/dev, we try to load a .env file if present.
func LoadAllConfigs(envFile string) {
	if envFile == "" {
		envFile = ".env"
	}

	// Decide if we should load the .env file based on common env flags
	env := strings.ToLower(strings.TrimSpace(os.Getenv("ENV")))
	if env == "" {
		env = strings.ToLower(strings.TrimSpace(os.Getenv("GO_ENV")))
	}
	if env == "" {
		env = strings.ToLower(strings.TrimSpace(os.Getenv("APP_ENV")))
	}

	// Only try to load .env for local/dev; ignore if missing
	if env == "" || env == "local" || env == "dev" || env == "development" {
		if err := godotenv.Load(envFile); err != nil {
			// Not fatal — in containers/prod there is no .env file by design
			log.Printf("can't load %s (continuing; using real env): %v", envFile, err)
		}
	}

	LoadApp()
	LoadDBCfg()
}

// FiberConfig func for configuration Fiber app.
func FiberConfig() fiber.Config {
	return fiber.Config{
		ReadTimeout: time.Second * time.Duration(AppCfg().ReadTimeout),
	}
}
