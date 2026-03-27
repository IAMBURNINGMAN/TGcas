package bot

import (
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"tgcasino/internal/payment"
	"tgcasino/internal/storage"
)

func (b *Bot) route(msg *tgbotapi.Message) {
	switch msg.Command() {
	case "start":
		b.handleStart(msg)
	case "balance":
		b.handleBalance(msg)
	case "promo":
		b.handlePromo(msg)
	case "dice":
		b.handleGame(msg, "dice")
	case "slots":
		b.handleGame(msg, "slots")
	}
}

func (b *Bot) handleStart(msg *tgbotapi.Message) {
	user, err := storage.GetOrCreate(b.db, msg.From.ID, msg.From.UserName)
	if err != nil {
		b.send(msg.Chat.ID, "Ой, что-то пошло не так 😢 Попробуй ещё раз!")
		return
	}
	b.send(msg.Chat.ID, fmt.Sprintf(
		"Привет, %s! 🎰 Добро пожаловать в казино Сабинки!\n\n"+
			"💰 Твой баланс: %d Сабинок\n\n"+
			"Что умеет Сабинка:\n"+
			"/balance — сколько Сабинок в кармане\n"+
			"/promo <код> — ввести промокод на Сабинки\n"+
			"/dice — сыграть в кубик 🎲\n"+
			"/slots — сыграть в слот 🎰",
		user.Username, user.Balance,
	))
}

func (b *Bot) handleBalance(msg *tgbotapi.Message) {
	balance, err := storage.GetBalance(b.db, msg.From.ID)
	if err != nil {
		b.send(msg.Chat.ID, "Сабинка тебя не знает 🤔 Напиши /start сначала!")
		return
	}
	b.send(msg.Chat.ID, fmt.Sprintf("💰 У тебя %d Сабинок!", balance))
}

func (b *Bot) handlePromo(msg *tgbotapi.Message) {
	code := msg.CommandArguments()
	if code == "" {
		b.send(msg.Chat.ID, "Напиши так: /promo <твой код>\nСабинка ждёт! 🎁")
		return
	}

	amount, err := payment.ApplyPromo(b.db, msg.From.ID, code)
	switch err {
	case nil:
		b.send(msg.Chat.ID, fmt.Sprintf("Промокод сработал! 🎉 +%d Сабинок на счёт!", amount))
	case payment.ErrPromoNotFound:
		b.send(msg.Chat.ID, "Сабинка такого промокода не знает 🙅")
	case payment.ErrPromoUsed:
		b.send(msg.Chat.ID, "Этот промокод уже использован 😔 Сабинки не резиновые!")
	default:
		b.send(msg.Chat.ID, "Ой, Сабинка сломалась 😢 Попробуй позже!")
	}
}

func (b *Bot) handleGame(msg *tgbotapi.Message, gameType string) {
	if _, err := storage.GetOrCreate(b.db, msg.From.ID, msg.From.UserName); err != nil {
		b.send(msg.Chat.ID, "Сабинка тебя не знает 🤔 Напиши /start сначала!")
		return
	}

	var text string
	switch gameType {
	case "dice":
		text = "На сколько Сабинок рискнёшь? 🎲"
	case "slots":
		text = "Сколько Сабинок ставим на удачу? 🎰"
	}

	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("10 Сабинок", gameType+":10"),
			tgbotapi.NewInlineKeyboardButtonData("50 Сабинок", gameType+":50"),
			tgbotapi.NewInlineKeyboardButtonData("100 Сабинок", gameType+":100"),
			tgbotapi.NewInlineKeyboardButtonData("500 Сабинок", gameType+":500"),
		),
	)
	reply := tgbotapi.NewMessage(msg.Chat.ID, text)
	reply.ReplyMarkup = keyboard
	b.api.Send(reply)
}
