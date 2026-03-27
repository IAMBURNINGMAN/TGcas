package main

import (
	"tgcasino/config"
	"tgcasino/internal/bot"
	"tgcasino/internal/storage"
)

func main() {
	cfg := config.Load()

	db := storage.Connect(cfg.DBUrl)
	defer db.Close()

	storage.RunMigrations(cfg.DBUrl)

	b := bot.New(cfg.BotToken, db)
	b.Start()
}
