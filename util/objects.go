package util

type presentType uint8
const (
	Coal presentType = iota
	Console
	Robot
	Doll
	Plush
	BoardGame
	Puzzle
	Book
	Lego
)

type Present struct {
	Type  presentType
}

type behaviourType uint8
const (
	Good behaviourType = iota
	Bad
)

type Address struct { // Note: Santa has an address of (0,0)
	X int
	Y int
}

type Child struct {
	Name      string
	Behaviour behaviourType
	Address   Address
	WishList  []Present
	Presents  []Present
}
