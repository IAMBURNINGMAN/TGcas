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
		b.showCabinetMsg(msg.Chat.ID, msg.From.ID, msg.From.UserName)
	case "promo":
		b.handlePromo(msg)
	// Обратная совместимость: команды игр открывают главное меню
	case "dice", "slots":
		b.handleStart(msg)
	}
}

func (b *Bot) handleStart(msg *tgbotapi.Message) {
	user, err := storage.GetOrCreate(b.db, msg.From.ID, msg.From.UserName)
	if err != nil {
		b.send(msg.Chat.ID, "Ой, что-то пошло не так 😢 Попробуй ещё раз!")
		return
	}
	reply := tgbotapi.NewMessage(msg.Chat.ID, welcomeText(user.Username, user.Balance))
	reply.ParseMode = "Markdown"
	reply.ReplyMarkup = MainMenuKeyboard()
	b.api.Send(reply)
}

func (b *Bot) handlePromo(msg *tgbotapi.Message) {
	code := msg.CommandArguments()
	if code == "" {
		b.send(msg.Chat.ID, "Напиши: /promo ТВОЙКОД\nСабинка ждёт! 🎁")
		return
	}
	amount, err := payment.ApplyPromo(b.db, msg.From.ID, code)
	switch err {
	case nil:
		b.send(msg.Chat.ID, fmt.Sprintf("🎉 Промокод сработал! +%d Сабинок на счёт!", amount))
	case payment.ErrPromoNotFound:
		b.send(msg.Chat.ID, "Сабинка такого промокода не знает 🙅")
	case payment.ErrPromoUsed:
		b.send(msg.Chat.ID, "Этот промокод уже использован 😔 Сабинки не резиновые!")
	default:
		b.send(msg.Chat.ID, "Ой, Сабинка сломалась 😢 Попробуй позже!")
	}
}

// showCabinetMsg — показать кабинет новым сообщением (из команды /balance).
func (b *Bot) showCabinetMsg(chatID int64, userID int64, username string) {
	balance, err := storage.GetBalance(b.db, userID)
	if err != nil {
		b.send(chatID, "Сабинка тебя не знает 🤔 Напиши /start сначала!")
		return
	}
	stats, err := storage.GetStats(b.db, userID)
	if err != nil {
		b.send(chatID, "Ошибка загрузки статистики 😢")
		return
	}
	msg := tgbotapi.NewMessage(chatID, cabinetText(username, balance, stats))
	msg.ParseMode = "Markdown"
	msg.ReplyMarkup = MainMenuKeyboard()
	b.api.Send(msg)
}

// showCabinet — показать кабинет редактированием сообщения (из меню).
func (b *Bot) showCabinet(chatID int64, msgID int, userID int64, username string) {
	balance, err := storage.GetBalance(b.db, userID)
	if err != nil {
		b.send(chatID, "Сабинка тебя не знает 🤔 Напиши /start сначала!")
		return
	}
	stats, err := storage.GetStats(b.db, userID)
	if err != nil {
		b.send(chatID, "Ошибка загрузки статистики 😢")
		return
	}
	b.editMD(chatID, msgID, cabinetText(username, balance, stats), MainMenuKeyboard())
}
