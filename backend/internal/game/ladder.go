package game

import "math/rand"

// LadderMultipliers — множители выплат по уровням лесенки (6 ступеней).
var LadderMultipliers = []float64{1.5, 2.0, 3.0, 5.0, 10.0, 20.0}

// LadderPayout — рассчитать выплату за текущий уровень.
func LadderPayout(bet int64, multIdx int) int64 {
	if multIdx < 0 || multIdx >= len(LadderMultipliers) {
		return 0
	}
	return int64(float64(bet) * LadderMultipliers[multIdx])
}

// LadderRoll — выжить на текущем уровне (60% шанс).
func LadderRoll() bool {
	return rand.Float64() < 0.60
}
