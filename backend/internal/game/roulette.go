package game

// RouletteResult — рулетка на основе 🎯 (значения 1-6).
// 1-3 = красное (×1.9), 4-5 = чёрное (×1.9), 6 = зелёное (×14).
func RouletteResult(bet int64, colour string, dartValue int) (won bool, payout int64) {
	var hit string
	switch dartValue {
	case 1, 2, 3:
		hit = "red"
	case 4, 5:
		hit = "black"
	case 6:
		hit = "green"
	}
	if colour != hit {
		return false, 0
	}
	switch colour {
	case "red", "black":
		return true, int64(float64(bet) * 1.9)
	case "green":
		return true, bet * 14
	}
	return false, 0
}
