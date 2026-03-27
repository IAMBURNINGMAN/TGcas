package game

// Значения 🎰 (1-64):
// 64 = джекпот 777 → x10, 22/43 = два совпадения → x2, остальное → проигрыш
func SlotsResult(bet int64, slotValue int) int64 {
	switch slotValue {
	case 64:
		return bet * 10
	case 22, 43:
		return bet * 2
	default:
		return 0
	}
}
