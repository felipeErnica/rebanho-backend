package util

func CalculatePercentageTrend(current float64, previous float64) float64 {
	if previous == 0 {
		return 0
	}

	if current == 0 {
		return 100
	}

	trend := ((current / previous) - 1) * 100
	return trend
}
