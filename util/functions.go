package util

func Check(e error) {
	if e != nil {
		panic(e)
	}
}

func Max(xs []int) int {
	max := 0
	for i, x := range xs {
		if i == 0 || x > max {
			max = x
		}
	}
	return max
}

func Min(xs []int) int {
	min := 0
	for x := range xs {
		if x == 0 || x < min {
			min = x
		}
	}
	return min
}