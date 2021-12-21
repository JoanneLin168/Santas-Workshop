package util

func Check(e error) {
	if e != nil {
		panic(e)
	}
}

// ChildInSlice - checks if a child is in the slice. Note: compares names - might need to make a comparison function for Child
func ChildInSlice(child Child, children []Child) bool {
	for c := range(children) {
		if child.Name == children[c].Name {
			return true
		}
	}
	return false
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