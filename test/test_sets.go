package test

import "workshop/util"

func GetSets() [][]util.Child {
	sets := [][]util.Child{
		SetA(), SetB(), SetC(), SetD(), SetE(),
	}
	return sets
}

// SetA - Test 0 children
func SetA() []util.Child {
	return []util.Child{}
}

// SetB - Test 1 child
func SetB() []util.Child {
	var children = []util.Child{
		{"Alice",
			util.Good,
			util.Address{X: 1, Y: 1},
			[]util.Present{{util.Doll}, {util.Book}, {util.Puzzle}, {util.Plush}},
			[]util.Present{},
		},
	}
	return children
}

// SetC - Test 3 children
func SetC() []util.Child {
	var children = []util.Child{
		{"Alice",
			util.Good,
			util.Address{X: 1, Y: 1},
			[]util.Present{{util.Doll}, {util.Book}, {util.Puzzle}, {util.Plush}},
			[]util.Present{},
		},
		{"Bob",
			util.Bad,
			util.Address{X: 0, Y: 5},
			[]util.Present{{util.Lego}, {util.Robot}, {util.Console}},
			[]util.Present{},
		},
		{"Charlie",
			util.Good,
			util.Address{X: -1, Y: -1},
			[]util.Present{{util.Book}, {util.BoardGame}, {util.Puzzle}, {util.Robot}, {util.Lego}},
			[]util.Present{},
		},
	}
	return children
}

// SetD - Test 8 children
func SetD() []util.Child {
	var children = []util.Child{
		{"Alice",
			util.Good,
			util.Address{X: 1, Y: 1},
			[]util.Present{{util.Doll}, {util.Book}, {util.Puzzle}, {util.Plush}},
			[]util.Present{},
		},
		{"Bob",
			util.Bad,
			util.Address{X: 0, Y: 5},
			[]util.Present{{util.Lego}, {util.Robot}, {util.Console}},
			[]util.Present{},
		},
		{"Charlie",
			util.Good,
			util.Address{X: -1, Y: -1},
			[]util.Present{{util.Book}, {util.BoardGame}, {util.Puzzle}, {util.Robot}, {util.Lego}},
			[]util.Present{},
		},
		{"David",
			util.Good,
			util.Address{X: 2, Y: 1},
			[]util.Present{{util.Doll}, {util.Book}, {util.Puzzle}, {util.Plush}},
			[]util.Present{},
		},
		{"Evie",
			util.Bad,
			util.Address{X: 0, Y: 1},
			[]util.Present{{util.Lego}, {util.Robot}, {util.Console}},
			[]util.Present{},
		},
		{"Fred",
			util.Bad,
			util.Address{X: 3, Y: 3},
			[]util.Present{{util.Book}, {util.BoardGame}, {util.Puzzle}, {util.Robot}, {util.Lego}},
			[]util.Present{},
		},
		{"Gemma",
			util.Bad,
			util.Address{X: 2, Y: 0},
			[]util.Present{{util.Doll}, {util.Book}, {util.Puzzle}, {util.Plush}},
			[]util.Present{},
		},
		{"Harry",
			util.Good,
			util.Address{X: 1, Y: 0},
			[]util.Present{{util.Lego}, {util.Robot}, {util.Console}},
			[]util.Present{},
		},
	}
	return children
}

// SetE - Test 10 children
func SetE() []util.Child {
	var children = []util.Child{
		{"Alice",
			util.Good,
			util.Address{X: 1, Y: 1},
			[]util.Present{{util.Doll}, {util.Book}, {util.Puzzle}, {util.Plush}},
			[]util.Present{},
		},
		{"Bob",
			util.Bad,
			util.Address{X: 0, Y: 5},
			[]util.Present{{util.Lego}, {util.Robot}, {util.Console}},
			[]util.Present{},
		},
		{"Charlie",
			util.Good,
			util.Address{X: -1, Y: -1},
			[]util.Present{{util.Book}, {util.BoardGame}, {util.Puzzle}, {util.Robot}, {util.Lego}},
			[]util.Present{},
		},
		{"David",
			util.Good,
			util.Address{X: 2, Y: 1},
			[]util.Present{{util.Doll}, {util.Book}, {util.Puzzle}, {util.Plush}},
			[]util.Present{},
		},
		{"Evie",
			util.Bad,
			util.Address{X: 0, Y: 1},
			[]util.Present{{util.Lego}, {util.Robot}, {util.Console}},
			[]util.Present{},
		},
		{"Fred",
			util.Bad,
			util.Address{X: 3, Y: 3},
			[]util.Present{{util.Book}, {util.BoardGame}, {util.Puzzle}, {util.Robot}, {util.Lego}},
			[]util.Present{},
		},
		{"Gemma",
			util.Bad,
			util.Address{X: 2, Y: 0},
			[]util.Present{{util.Doll}, {util.Book}, {util.Puzzle}, {util.Plush}},
			[]util.Present{},
		},
		{"Harry",
			util.Good,
			util.Address{X: 1, Y: 0},
			[]util.Present{{util.Lego}, {util.Robot}, {util.Console}},
			[]util.Present{},
		},
		{"Isabelle",
			util.Good,
			util.Address{X: -3, Y: 1},
			[]util.Present{{util.Book}, {util.BoardGame}, {util.Puzzle}, {util.Robot}, {util.Lego}},
			[]util.Present{},
		},
		{"Jake",
			util.Bad,
			util.Address{X: -2, Y: -4},
			[]util.Present{{util.Doll}, {util.Book}, {util.Puzzle}, {util.Plush}},
			[]util.Present{},
		},
	}
	return children
}