package game

// CoinflipResult — монетка на основе 🎲 (значения 1-6).
// 1-3 = орёл, 4-6 = решка. Выплата ×1.9 от ставки.
func CoinflipResult(bet int64, choice string, diceValue int) (won bool, payout int64) {
	result := "heads"
	if diceValue > 3 {
		result = "tails"
	}
	if choice == result {
		return true, int64(float64(bet) * 1.9)
	}
	return false, 0
}
