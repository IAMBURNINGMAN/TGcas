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

// handleCrashBet — списать ставку и показать выбор целевого множителя.
func (b *Bot) handleCrashBet(cb *tgbotapi.CallbackQuery, betStr string) {
	bet, err := strconv.ParseInt(betStr, 10, 64)
	if err != nil {
		return
	}
	if err := wallet.Debit(b.db, cb.From.ID, bet, "bet:crash"); err != nil {
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
			"💣 *Краш!* Ставка: *%d Сабинок*\n\n"+
				"🚀 Ракета взлетает...\n"+
				"Выбери цель — если краш случится ДО неё, ты проиграешь!\n\n"+
				"_Выше цель = выше риск = выше выигрыш_ 🔥",
			bet,
		),
		CrashChoiceKeyboard(bet),
	)
}

// handleCrashGo — рассчитать краш и показать результат.
// rest: "bet:targetX100" (множитель ×100, например 200 = ×2.0)
func (b *Bot) handleCrashGo(cb *tgbotapi.CallbackQuery, rest string) {
	parts := strings.SplitN(rest, ":", 2)
	if len(parts) != 2 {
		return
	}
	bet, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil {
		return
	}
	targetX100, err := strconv.ParseInt(parts[1], 10, 64)
	if err != nil {
		return
	}
	targetMult := float64(targetX100) / 100.0
	b.api.Request(tgbotapi.NewCallback(cb.ID, "🚀"))

	crashed, actualMult, payout := game.CrashResult(bet, targetMult)
	wallet.RecordResult(b.db, cb.From.ID, bet, payout, "win:crash")
	balance, _ := storage.GetBalance(b.db, cb.From.ID)

	var text string
	if !crashed {
		text = fmt.Sprintf(
			"💣 Краш на *×%.2f*\n\n"+
				"🎯 Твоя цель *×%.1f* — достигнута!\n\n"+
				"✅ Выигрыш: *+%d Сабинок*\n"+
				"💰 Баланс: *%d*",
			actualMult, targetMult, payout, balance,
		)
	} else {
		text = fmt.Sprintf(
			"💣 Краш на *×%.2f*\n\n"+
				"🎯 Твоя цель *×%.1f* — не достигнута!\n\n"+
				"❌ Потеряно: *%d Сабинок*\n"+
				"💰 Баланс: *%d*\n\n"+
				"_Риск — дело благородное! Попробуй ещё!_",
			actualMult, targetMult, bet, balance,
		)
	}
	b.sendMD(cb.Message.Chat.ID, text, PlayAgainKeyboard("crash"))
}
