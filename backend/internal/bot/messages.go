package bot

import (
	"fmt"

	"tgcasino/internal/storage"
)

// welcomeText — приветственное сообщение главного меню.
func welcomeText(username string, balance int64) string {
	return fmt.Sprintf(
		"🎰 Добро пожаловать в казино *Сабинки*, %s!\n\n"+
			"💰 Твой баланс: *%d Сабинок*\n\n"+
			"🎮 Выбери игру и испытай удачу!\n"+
			"🎁 Есть промокод? Нажми кнопку ниже!\n\n"+
			"_Чем выше риск — тем больше выигрыш_ 🔥",
		username, balance,
	)
}

// cabinetText — личный кабинет с подробной статистикой.
func cabinetText(username string, balance int64, s *storage.UserStats) string {
	var winRate float64
	if s.GamesPlayed > 0 {
		winRate = float64(s.Wins) / float64(s.GamesPlayed) * 100
	}
	return fmt.Sprintf(
		"👤 *Личный кабинет* @%s\n\n"+
			"💰 Баланс: *%d Сабинок*\n\n"+
			"📊 *Статистика:*\n"+
			"├ Игр сыграно: *%d*\n"+
			"├ Побед: *%d* (%.1f%%)\n"+
			"├ Поставлено всего: *%d*\n"+
			"└ Выиграно всего: *%d*\n\n"+
			"🏆 Лучший выигрыш: *%d Сабинок*\n"+
			"📅 С нами с: %s",
		username, balance,
		s.GamesPlayed, s.Wins, winRate,
		s.TotalWagered, s.TotalWon,
		s.BestWin, s.MemberSince,
	)
}
