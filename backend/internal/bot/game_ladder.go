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

// handleLadderStart — списать ставку и показать первый уровень.
func (b *Bot) handleLadderStart(cb *tgbotapi.CallbackQuery, betStr string) {
	bet, err := strconv.ParseInt(betStr, 10, 64)
	if err != nil {
		return
	}
	if err := wallet.Debit(b.db, cb.From.ID, bet, "bet:ladder"); err != nil {
		if err == wallet.ErrInsufficientFunds {
			b.api.Request(tgbotapi.NewCallback(cb.ID, "Сабинок не хватает! 😢"))
		} else {
			b.api.Request(tgbotapi.NewCallback(cb.ID, "Ошибка 😢"))
		}
		return
	}
	b.api.Request(tgbotapi.NewCallback(cb.ID, ""))
	b.editMD(cb.Message.Chat.ID, cb.Message.MessageID,
		ladderLevelText(bet, 0),
		LadderActionKeyboard(bet, 0),
	)
}

// handleLadderTake — забрать текущий выигрыш.
// rest: "bet:multIdx"
func (b *Bot) handleLadderTake(cb *tgbotapi.CallbackQuery, rest string) {
	bet, multIdx, err := parseLadderRest(rest)
	if err != nil {
		return
	}
	b.api.Request(tgbotapi.NewCallback(cb.ID, ""))

	payout := game.LadderPayout(bet, multIdx)
	wallet.RecordResult(b.db, cb.From.ID, bet, payout, "win:ladder")
	balance, _ := storage.GetBalance(b.db, cb.From.ID)

	text := fmt.Sprintf(
		"🪜 *Умно! Забираешь выигрыш!*\n\n"+
			"✅ Уровень %d — *×%.1f*\n"+
			"💰 Получено: *+%d Сабинок*\n"+
			"💰 Баланс: *%d*",
		multIdx+1, game.LadderMultipliers[multIdx], payout, balance,
	)
	b.editMD(cb.Message.Chat.ID, cb.Message.MessageID, text, PlayAgainKeyboard("ladder"))
}

// handleLadderRisk — рискнуть на следующем уровне.
// rest: "bet:multIdx"
func (b *Bot) handleLadderRisk(cb *tgbotapi.CallbackQuery, rest string) {
	bet, multIdx, err := parseLadderRest(rest)
	if err != nil {
		return
	}
	b.api.Request(tgbotapi.NewCallback(cb.ID, "🎲 Рулим..."))

	if game.LadderRoll() {
		nextIdx := multIdx + 1
		if nextIdx >= len(game.LadderMultipliers) {
			nextIdx = len(game.LadderMultipliers) - 1
		}
		b.editMD(cb.Message.Chat.ID, cb.Message.MessageID,
			ladderLevelText(bet, nextIdx),
			LadderActionKeyboard(bet, nextIdx),
		)
		return
	}

	// Проигрыш
	wallet.RecordResult(b.db, cb.From.ID, bet, 0, "loss:ladder")
	balance, _ := storage.GetBalance(b.db, cb.From.ID)
	text := fmt.Sprintf(
		"💥 *Лесенка сломалась!*\n\n"+
			"Ты дошёл до уровня *%d* (×%.1f), но упал...\n\n"+
			"❌ Потеряно: *%d Сабинок*\n"+
			"💰 Баланс: *%d*\n\n"+
			"_Удача не за горами! Попробуй снова!_",
		multIdx+1, game.LadderMultipliers[multIdx], bet, balance,
	)
	b.editMD(cb.Message.Chat.ID, cb.Message.MessageID, text, PlayAgainKeyboard("ladder"))
}

// ladderLevelText — текст для текущего уровня лесенки.
func ladderLevelText(bet int64, multIdx int) string {
	mult := game.LadderMultipliers[multIdx]
	payout := game.LadderPayout(bet, multIdx)
	level := multIdx + 1
	total := len(game.LadderMultipliers)

	var header string
	switch {
	case multIdx == 0:
		header = "🪜 *Лесенка запущена!*"
	case multIdx >= total-1:
		header = "👑 *МАКСИМАЛЬНЫЙ УРОВЕНЬ — ДЖЕКПОТ!*"
	default:
		header = "✅ *Уровень пройден! Идёшь выше!*"
	}

	text := fmt.Sprintf(
		"%s\n\n"+
			"Ставка: *%d Сабинок*\n"+
			"━━━━━━━━━━\n"+
			"📍 Уровень *%d из %d*\n"+
			"✨ Текущий выигрыш: *%d* (×%.1f)\n"+
			"━━━━━━━━━━",
		header, bet, level, total, payout, mult,
	)

	if multIdx < total-1 {
		nextMult := game.LadderMultipliers[multIdx+1]
		nextPayout := game.LadderPayout(bet, multIdx+1)
		text += fmt.Sprintf(
			"\n🎯 Следующий: *×%.1f* = %d Сабинок\n"+
				"⚠️ Шанс пройти: *60%%*",
			nextMult, nextPayout,
		)
	}

	return text
}

// parseLadderRest — распарсить "bet:multIdx" из callback rest.
func parseLadderRest(rest string) (bet int64, multIdx int, err error) {
	parts := strings.SplitN(rest, ":", 2)
	if len(parts) != 2 {
		return 0, 0, fmt.Errorf("invalid ladder rest")
	}
	bet, err = strconv.ParseInt(parts[0], 10, 64)
	if err != nil {
		return
	}
	idx, err := strconv.ParseInt(parts[1], 10, 64)
	if err != nil {
		return
	}
	multIdx = int(idx)
	return
}
