package bot

import (
	"fmt"
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"tgcasino/internal/game"
	"tgcasino/internal/wallet"
)

func (b *Bot) handleCallback(cb *tgbotapi.CallbackQuery) {
	parts := strings.SplitN(cb.Data, ":", 2)
	if len(parts) != 2 {
		return
	}
	gameType, betStr := parts[0], parts[1]

	bet, err := strconv.ParseInt(betStr, 10, 64)
	if err != nil {
		return
	}

	if err := wallet.Debit(b.db, cb.From.ID, bet, "bet:"+gameType); err != nil {
		if err == wallet.ErrInsufficientFunds {
			b.api.Request(tgbotapi.NewCallback(cb.ID, "Сабинок не хватает! 😢"))
		} else {
			b.api.Request(tgbotapi.NewCallback(cb.ID, "Ошибка 😢"))
		}
		return
	}
	b.api.Request(tgbotapi.NewCallback(cb.ID, ""))

	var emoji string
	switch gameType {
	case "dice":
		emoji = "🎲"
	case "slots":
		emoji = "🎰"
	}

	diceMsg := tgbotapi.NewDice(cb.Message.Chat.ID)
	diceMsg.Emoji = emoji
	sent, err := b.api.Send(diceMsg)
	if err != nil || sent.Dice == nil {
		wallet.Credit(b.db, cb.From.ID, bet, "refund:"+gameType)
		b.send(cb.Message.Chat.ID, "Сабинка не смогла бросить кубик 😢 Ставка возвращена!")
		return
	}

	var winAmount int64
	switch gameType {
	case "dice":
		winAmount = game.DiceResult(bet, sent.Dice.Value)
	case "slots":
		winAmount = game.SlotsResult(bet, sent.Dice.Value)
	}

	if winAmount > 0 {
		wallet.Credit(b.db, cb.From.ID, winAmount, "win:"+gameType)
		b.send(cb.Message.Chat.ID, fmt.Sprintf("Сабинка улыбается тебе! 🎉 +%d Сабинок!", winAmount))
	} else {
		b.send(cb.Message.Chat.ID, fmt.Sprintf("Не повезло 😔 %d Сабинок сгорело. Сабинка верит в тебя — пробуй снова!", bet))
	}
}
