package base

func Lerp(a float64, b float64, f float64) float64 {
	return (a * (1.0 - f)) + (b * f)
}
