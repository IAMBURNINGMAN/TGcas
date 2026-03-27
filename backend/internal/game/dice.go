package game

// Множители для 🎲 (значения 1-6):
// 1-2 → проигрыш, 3 → возврат ставки, 4 → x2, 5 → x3, 6 → x5
var diceMultipliers = map[int]int64{
	1: 0,
	2: 0,
	3: 1,
	4: 2,
	5: 3,
	6: 5,
}

func DiceResult(bet int64, diceValue int) int64 {
	mult, ok := diceMultipliers[diceValue]
	if !ok {
		return 0
	}
	return bet * mult
}
