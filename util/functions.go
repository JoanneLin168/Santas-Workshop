package util

func Check(e error) {
	if e != nil {
		panic(e)
	}
}

func RemoveChild(i int, xs []Child) []Child {
	return append(xs[:i], xs[i+1:]...)
}

func MaxInt(xs []int) int {
	max := xs[0]
	for i := range xs {
		x := xs[i]
		if x > max {
			max = x
		}
	}
	return max
}

func MinInt(xs []int) int {
	min := xs[0]
	for i := range xs {
		x := xs[i]
		if x < min {
			min = x
		}
	}
	return min
}

func MaxFloat64(xs []float64) float64 {
	max := xs[0]
	for i := range xs {
		x := xs[i]
		if x > max {
			max = x
		}
	}
	return max
}

func MinFloat64(xs []float64) float64 {
	min := xs[0]
	for i := range xs {
		x := xs[i]
		if x < min {
			min = x
		}
	}
	return min
}