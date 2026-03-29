package bot

import (
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"tgcasino/internal/game"
)

// MainMenuKeyboard — главное меню казино.
func MainMenuKeyboard() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🎲 Кубик", "menu:dice"),
			tgbotapi.NewInlineKeyboardButtonData("🎰 Слоты", "menu:slots"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🪙 Монетка", "menu:coin"),
			tgbotapi.NewInlineKeyboardButtonData("🎯 Рулетка", "menu:roulette"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🪜 Лесенка", "menu:ladder"),
			tgbotapi.NewInlineKeyboardButtonData("💣 Краш", "menu:crash"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("👤 Кабинет", "menu:cabinet"),
			tgbotapi.NewInlineKeyboardButtonData("🎁 Промокод", "menu:promo"),
		),
	)
}

// BetKeyboard — выбор ставки для игры.
func BetKeyboard(gameType string) tgbotapi.InlineKeyboardMarkup {
	btn := func(label, amt string) tgbotapi.InlineKeyboardButton {
		return tgbotapi.NewInlineKeyboardButtonData(label, gameType+":"+amt)
	}
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			btn("10 💰", "10"), btn("50 💰", "50"),
			btn("100 💰", "100"), btn("500 💰", "500"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("◀ Меню", "menu:main"),
		),
	)
}

// CoinChoiceKeyboard — орёл или решка после ставки.
func CoinChoiceKeyboard(bet int64) tgbotapi.InlineKeyboardMarkup {
	s := fmt.Sprintf("%d", bet)
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🦅 Орёл", "coin_flip:"+s+":heads"),
			tgbotapi.NewInlineKeyboardButtonData("🌐 Решка", "coin_flip:"+s+":tails"),
		),
	)
}

// RouletteChoiceKeyboard — выбор цвета после ставки.
func RouletteChoiceKeyboard(bet int64) tgbotapi.InlineKeyboardMarkup {
	s := fmt.Sprintf("%d", bet)
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🔴 Красное ×1.9", "roulette_spin:"+s+":red"),
			tgbotapi.NewInlineKeyboardButtonData("⚫ Чёрное ×1.9", "roulette_spin:"+s+":black"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🟢 Зелёное ×14 🔥", "roulette_spin:"+s+":green"),
		),
	)
}

// CrashChoiceKeyboard — выбор целевого множителя.
// Множитель кодируется как int ×100 (например 200 = ×2.0).
func CrashChoiceKeyboard(bet int64) tgbotapi.InlineKeyboardMarkup {
	s := fmt.Sprintf("%d", bet)
	m := game.CrashMultipliers
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(fmt.Sprintf("×%.1f", m[0]), "crash_go:"+s+":150"),
			tgbotapi.NewInlineKeyboardButtonData(fmt.Sprintf("×%.1f", m[1]), "crash_go:"+s+":200"),
			tgbotapi.NewInlineKeyboardButtonData(fmt.Sprintf("×%.1f", m[2]), "crash_go:"+s+":300"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(fmt.Sprintf("×%.1f", m[3]), "crash_go:"+s+":500"),
			tgbotapi.NewInlineKeyboardButtonData(fmt.Sprintf("×%.1f 🔥", m[4]), "crash_go:"+s+":1000"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("◀ Меню", "menu:main"),
		),
	)
}

// LadderActionKeyboard — забрать или рискнуть.
// На последнем уровне только кнопка "Забрать".
func LadderActionKeyboard(bet int64, multIdx int) tgbotapi.InlineKeyboardMarkup {
	payout := game.LadderPayout(bet, multIdx)
	bs := fmt.Sprintf("%d", bet)
	ms := fmt.Sprintf("%d", multIdx)
	takeBtn := tgbotapi.NewInlineKeyboardButtonData(
		fmt.Sprintf("💰 Забрать %d", payout), "ladder_take:"+bs+":"+ms,
	)
	if multIdx >= len(game.LadderMultipliers)-1 {
		return tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(takeBtn),
		)
	}
	riskBtn := tgbotapi.NewInlineKeyboardButtonData("🔥 Рискнуть!", "ladder_risk:"+bs+":"+ms)
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(takeBtn, riskBtn),
	)
}

// PlayAgainKeyboard — кнопки после завершения игры.
func PlayAgainKeyboard(gameType string) tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🔄 Ещё раз", "menu:"+gameType),
			tgbotapi.NewInlineKeyboardButtonData("◀ Меню", "menu:main"),
		),
	)
}
