package util

func CalculatePercentageTrend(current float64, previous float64) float64 {
	return (current / previous) - 1
}
