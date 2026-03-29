package game

import "math/rand"

// CrashMultipliers — доступные целевые множители для краша.
var CrashMultipliers = []float64{1.5, 2.0, 3.0, 5.0, 10.0}

// CrashResult — краш: выиграть если случайный множитель >= цели.
// Распределение 1/(1-r*0.95) — чаще низкие значения, даёт перевес казино.
func CrashResult(bet int64, targetMult float64) (crashed bool, actualMult float64, payout int64) {
	r := rand.Float64()
	actual := 1.0 / (1.0 - r*0.95)
	if actual > 20.0 {
		actual = 20.0
	}
	actual = float64(int64(actual*100)) / 100.0

	if actual >= targetMult {
		return false, actual, int64(float64(bet) * targetMult)
	}
	return true, actual, 0
}
