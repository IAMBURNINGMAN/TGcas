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

// handleCallback — главный роутер callback-кнопок.
func (b *Bot) handleCallback(cb *tgbotapi.CallbackQuery) {
	parts := strings.SplitN(cb.Data, ":", 2)
	action := parts[0]
	rest := ""
	if len(parts) == 2 {
		rest = parts[1]
	}

	switch action {
	case "menu":
		b.handleMenuCB(cb, rest)
	case "dice":
		b.handleDiceCB(cb, rest)
	case "slots":
		b.handleSlotsCB(cb, rest)
	case "coin":
		b.handleCoinBet(cb, rest)
	case "coin_flip":
		b.handleCoinFlip(cb, rest)
	case "roulette":
		b.handleRouletteBet(cb, rest)
	case "roulette_spin":
		b.handleRouletteSpin(cb, rest)
	case "crash":
		b.handleCrashBet(cb, rest)
	case "crash_go":
		b.handleCrashGo(cb, rest)
	case "ladder":
		b.handleLadderStart(cb, rest)
	case "ladder_take":
		b.handleLadderTake(cb, rest)
	case "ladder_risk":
		b.handleLadderRisk(cb, rest)
	}
}

// handleMenuCB — навигация по меню.
func (b *Bot) handleMenuCB(cb *tgbotapi.CallbackQuery, action string) {
	b.api.Request(tgbotapi.NewCallback(cb.ID, ""))
	chatID := cb.Message.Chat.ID
	msgID := cb.Message.MessageID

	switch action {
	case "main":
		balance, _ := storage.GetBalance(b.db, cb.From.ID)
		b.editMD(chatID, msgID, welcomeText(cb.From.UserName, balance), MainMenuKeyboard())
	case "cabinet":
		b.showCabinet(chatID, msgID, cb.From.ID, cb.From.UserName)
	case "promo":
		b.editMD(chatID, msgID,
			"🎁 *Промокод*\n\nНапиши команду в чат:\n`/promo ТВОЙКОД`\n\nИ Сабинки упадут на счёт! 🤑",
			tgbotapi.NewInlineKeyboardMarkup(
				tgbotapi.NewInlineKeyboardRow(
					tgbotapi.NewInlineKeyboardButtonData("◀ Назад в меню", "menu:main"),
				),
			),
		)
	case "dice":
		b.editMD(chatID, msgID,
			"🎲 *Кубик!*\n\n1-2 = проигрыш, 3 = возврат, 4 = ×2, 5 = ×3, 6 = ×5 🔥\n\nВыбери ставку:",
			BetKeyboard("dice"),
		)
	case "slots":
		b.editMD(chatID, msgID,
			"🎰 *Слоты!*\n\n777 = ×10 джекпот 🔥, 2 совпадения = ×2\n\nВыбери ставку:",
			BetKeyboard("slots"),
		)
	case "coin":
		b.editMD(chatID, msgID,
			"🪙 *Монетка!*\n\nОрёл или решка — шанс 50/50, выплата ×1.9\n\nВыбери ставку:",
			BetKeyboard("coin"),
		)
	case "roulette":
		b.editMD(chatID, msgID,
			"🎯 *Рулетка!*\n\nКрасное/чёрное ×1.9, зелёное ×14 🔥\n\nВыбери ставку:",
			BetKeyboard("roulette"),
		)
	case "ladder":
		b.editMD(chatID, msgID,
			"🪜 *Лесенка!*\n\n×1.5 → ×2 → ×3 → ×5 → ×10 → ×20\n"+
				"На каждом уровне: забрать или рискнуть выше!\n"+
				"Шанс пройти каждый уровень: 60%\n\nВыбери ставку:",
			BetKeyboard("ladder"),
		)
	case "crash":
		b.editMD(chatID, msgID,
			"💣 *Краш!*\n\nВыбери цель и молись, чтобы ракета не взорвалась раньше! 🚀\n\nВыбери ставку:",
			BetKeyboard("crash"),
		)
	}
}

// handleDiceCB — бросить кубик.
func (b *Bot) handleDiceCB(cb *tgbotapi.CallbackQuery, betStr string) {
	bet, err := strconv.ParseInt(betStr, 10, 64)
	if err != nil {
		return
	}
	if err := wallet.Debit(b.db, cb.From.ID, bet, "bet:dice"); err != nil {
		if err == wallet.ErrInsufficientFunds {
			b.api.Request(tgbotapi.NewCallback(cb.ID, "Сабинок не хватает! 😢"))
		} else {
			b.api.Request(tgbotapi.NewCallback(cb.ID, "Ошибка 😢"))
		}
		return
	}
	b.api.Request(tgbotapi.NewCallback(cb.ID, ""))

	diceMsg := tgbotapi.NewDice(cb.Message.Chat.ID)
	diceMsg.Emoji = "🎲"
	sent, err := b.api.Send(diceMsg)
	if err != nil || sent.Dice == nil {
		wallet.Credit(b.db, cb.From.ID, bet, "refund:dice")
		b.send(cb.Message.Chat.ID, "Кубик не бросился 😢 Ставка возвращена!")
		return
	}

	winAmount := game.DiceResult(bet, sent.Dice.Value)
	wallet.RecordResult(b.db, cb.From.ID, bet, winAmount, "win:dice")
	balance, _ := storage.GetBalance(b.db, cb.From.ID)

	var text string
	if winAmount > 0 {
		text = fmt.Sprintf(
			"🎲 Выпало *%d*!\n\n✅ Выигрыш: *+%d Сабинок*\n💰 Баланс: *%d*",
			sent.Dice.Value, winAmount, balance,
		)
	} else {
		text = fmt.Sprintf(
			"🎲 Выпало *%d* — мимо!\n\n❌ Потеряно: *%d Сабинок*\n💰 Баланс: *%d*\n\n_Попробуй ещё!_",
			sent.Dice.Value, bet, balance,
		)
	}
	b.sendMD(cb.Message.Chat.ID, text, PlayAgainKeyboard("dice"))
}

// handleSlotsCB — крутануть слоты.
func (b *Bot) handleSlotsCB(cb *tgbotapi.CallbackQuery, betStr string) {
	bet, err := strconv.ParseInt(betStr, 10, 64)
	if err != nil {
		return
	}
	if err := wallet.Debit(b.db, cb.From.ID, bet, "bet:slots"); err != nil {
		if err == wallet.ErrInsufficientFunds {
			b.api.Request(tgbotapi.NewCallback(cb.ID, "Сабинок не хватает! 😢"))
		} else {
			b.api.Request(tgbotapi.NewCallback(cb.ID, "Ошибка 😢"))
		}
		return
	}
	b.api.Request(tgbotapi.NewCallback(cb.ID, ""))

	diceMsg := tgbotapi.NewDice(cb.Message.Chat.ID)
	diceMsg.Emoji = "🎰"
	sent, err := b.api.Send(diceMsg)
	if err != nil || sent.Dice == nil {
		wallet.Credit(b.db, cb.From.ID, bet, "refund:slots")
		b.send(cb.Message.Chat.ID, "Слот завис 😢 Ставка возвращена!")
		return
	}

	winAmount := game.SlotsResult(bet, sent.Dice.Value)
	wallet.RecordResult(b.db, cb.From.ID, bet, winAmount, "win:slots")
	balance, _ := storage.GetBalance(b.db, cb.From.ID)

	var text string
	switch {
	case winAmount >= bet*10:
		text = fmt.Sprintf(
			"🎰 *777 — ДЖЕКПОТ!* 🎉🎉🎉\n\n✅ Выигрыш: *+%d Сабинок* (×10)\n💰 Баланс: *%d*",
			winAmount, balance,
		)
	case winAmount > 0:
		text = fmt.Sprintf(
			"🎰 Два совпадения!\n\n✅ Выигрыш: *+%d Сабинок* (×2)\n💰 Баланс: *%d*",
			winAmount, balance,
		)
	default:
		text = fmt.Sprintf(
			"🎰 Нет совпадений...\n\n❌ Потеряно: *%d Сабинок*\n💰 Баланс: *%d*\n\n_Джекпот уже близко!_",
			bet, balance,
		)
	}
	b.sendMD(cb.Message.Chat.ID, text, PlayAgainKeyboard("slots"))
}
