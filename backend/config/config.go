package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	BotToken string
	DBUrl    string
}

func Load() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("файл .env не найден, читаем из окружения")
	}

	return &Config{
		BotToken: mustGet("BOT_TOKEN"),
		DBUrl:    mustGet("DATABASE_URL"),
	}
}

func mustGet(key string) string {
	v := os.Getenv(key)
	if v == "" {
		log.Fatalf("переменная окружения %s не задана", key)
	}
	return v
}
