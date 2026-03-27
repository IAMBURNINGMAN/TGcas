package bot

import (
	"database/sql"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Bot struct {
	api *tgbotapi.BotAPI
	db  *sql.DB
}

func New(token string, db *sql.DB) *Bot {
	api, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Fatalf("tgbotapi.NewBotAPI: %v", err)
	}
	log.Printf("бот запущен: @%s", api.Self.UserName)
	return &Bot{api: api, db: db}
}

func (b *Bot) Start() {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := b.api.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil {
			b.route(update.Message)
		} else if update.CallbackQuery != nil {
			b.handleCallback(update.CallbackQuery)
		}
	}
}

func (b *Bot) send(chatID int64, text string) {
	msg := tgbotapi.NewMessage(chatID, text)
	b.api.Send(msg)
}
