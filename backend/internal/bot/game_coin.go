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

// handleCoinBet — списать ставку и показать выбор орёл/решка.
func (b *Bot) handleCoinBet(cb *tgbotapi.CallbackQuery, betStr string) {
	bet, err := strconv.ParseInt(betStr, 10, 64)
	if err != nil {
		return
	}
	if err := wallet.Debit(b.db, cb.From.ID, bet, "bet:coin"); err != nil {
		if err == wallet.ErrInsufficientFunds {
			b.api.Request(tgbotapi.NewCallback(cb.ID, "Сабинок не хватает! 😢"))
		} else {
			b.api.Request(tgbotapi.NewCallback(cb.ID, "Ошибка 😢"))
		}
		return
	}
	b.api.Request(tgbotapi.NewCallback(cb.ID, ""))
	b.editMD(cb.Message.Chat.ID, cb.Message.MessageID,
		fmt.Sprintf("🪙 Ставка *%d Сабинок* принята!\n\nОрёл или решка?", bet),
		CoinChoiceKeyboard(bet),
	)
}

// handleCoinFlip — бросить монетку и показать результат.
// rest: "bet:choice"
func (b *Bot) handleCoinFlip(cb *tgbotapi.CallbackQuery, rest string) {
	parts := strings.SplitN(rest, ":", 2)
	if len(parts) != 2 {
		return
	}
	bet, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil {
		return
	}
	choice := parts[1]
	b.api.Request(tgbotapi.NewCallback(cb.ID, ""))

	diceMsg := tgbotapi.NewDice(cb.Message.Chat.ID)
	diceMsg.Emoji = "🎲"
	sent, err := b.api.Send(diceMsg)
	if err != nil || sent.Dice == nil {
		wallet.Credit(b.db, cb.From.ID, bet, "refund:coin")
		b.send(cb.Message.Chat.ID, "Монетка потерялась 😢 Ставка возвращена!")
		return
	}

	won, payout := game.CoinflipResult(bet, choice, sent.Dice.Value)
	wallet.RecordResult(b.db, cb.From.ID, bet, payout, "win:coin")
	balance, _ := storage.GetBalance(b.db, cb.From.ID)

	choiceLabel := "🌐 Решка"
	if choice == "heads" {
		choiceLabel = "🦅 Орёл"
	}

	var text string
	if won {
		text = fmt.Sprintf(
			"🪙 Ты поставил на %s!\n\n"+
				"✅ Угадал! *+%d Сабинок*\n"+
				"💰 Баланс: *%d*",
			choiceLabel, payout, balance,
		)
	} else {
		text = fmt.Sprintf(
			"🪙 Ты поставил на %s...\n\n"+
				"❌ Не угадал! *-%d Сабинок*\n"+
				"💰 Баланс: *%d*\n\n"+
				"_Повезёт в следующий раз!_",
			choiceLabel, bet, balance,
		)
	}
	b.sendMD(cb.Message.Chat.ID, text, PlayAgainKeyboard("coin"))
}
