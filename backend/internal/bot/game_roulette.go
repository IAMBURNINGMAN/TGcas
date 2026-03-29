package bot

import (
	"fmt"
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"tgcasino/internal/game"
	"tgcasino/internal/storage"
	"tgcasino/internal/wallet"
)

// handleRouletteBet — списать ставку и показать выбор цвета.
func (b *Bot) handleRouletteBet(cb *tgbotapi.CallbackQuery, betStr string) {
	bet, err := strconv.ParseInt(betStr, 10, 64)
	if err != nil {
		return
	}
	if err := wallet.Debit(b.db, cb.From.ID, bet, "bet:roulette"); err != nil {
		if err == wallet.ErrInsufficientFunds {
			b.api.Request(tgbotapi.NewCallback(cb.ID, "Сабинок не хватает! 😢"))
		} else {
			b.api.Request(tgbotapi.NewCallback(cb.ID, "Ошибка 😢"))
		}
		return
	}
	b.api.Request(tgbotapi.NewCallback(cb.ID, ""))
	b.editMD(cb.Message.Chat.ID, cb.Message.MessageID,
		fmt.Sprintf(
			"🎯 *Рулетка!* Ставка: *%d Сабинок*\n\n"+
				"🔴 Красное — шанс 50%%, выплата ×1.9\n"+
				"⚫ Чёрное — шанс 33%%, выплата ×1.9\n"+
				"🟢 Зелёное — шанс 17%%, выплата ×14 🔥\n\n"+
				"Выбирай цвет!",
			bet,
		),
		RouletteChoiceKeyboard(bet),
	)
}

// handleRouletteSpin — крутануть рулетку и показать результат.
// rest: "bet:colour"
func (b *Bot) handleRouletteSpin(cb *tgbotapi.CallbackQuery, rest string) {
	parts := strings.SplitN(rest, ":", 2)
	if len(parts) != 2 {
		return
	}
	bet, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil {
		return
	}
	colour := parts[1]
	b.api.Request(tgbotapi.NewCallback(cb.ID, ""))

	diceMsg := tgbotapi.NewDice(cb.Message.Chat.ID)
	diceMsg.Emoji = "🎯"
	sent, err := b.api.Send(diceMsg)
	if err != nil || sent.Dice == nil {
		wallet.Credit(b.db, cb.From.ID, bet, "refund:roulette")
		b.send(cb.Message.Chat.ID, "Рулетка сломалась 😢 Ставка возвращена!")
		return
	}

	won, payout := game.RouletteResult(bet, colour, sent.Dice.Value)
	wallet.RecordResult(b.db, cb.From.ID, bet, payout, "win:roulette")
	balance, _ := storage.GetBalance(b.db, cb.From.ID)

	labels := map[string]string{"red": "🔴 Красное", "black": "⚫ Чёрное", "green": "🟢 Зелёное"}
	label := labels[colour]

	var text string
	if won {
		prefix := ""
		if colour == "green" {
			prefix = "🎉🎉🎉 *ЗЕЛЁНОЕ — ДЖЕКПОТ!* 🎉🎉🎉\n\n"
		}
		text = prefix + fmt.Sprintf(
			"🎯 %s — выпало!\n\n"+
				"✅ Выигрыш! *+%d Сабинок*\n"+
				"💰 Баланс: *%d*",
			label, payout, balance,
		)
	} else {
		text = fmt.Sprintf(
			"🎯 Ты поставил на %s, но не повезло...\n\n"+
				"❌ *-%d Сабинок*\n"+
				"💰 Баланс: *%d*\n\n"+
				"_Удача на следующем круге!_",
			label, bet, balance,
		)
	}
	b.sendMD(cb.Message.Chat.ID, text, PlayAgainKeyboard("roulette"))
}
