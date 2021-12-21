package util

type PresentType uint8
const (
	Coal PresentType = iota
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
	Type PresentType
}

type BehaviourType uint8
const (
	Good BehaviourType = iota
	Bad
)

type Address struct { // Note: Santa has an address of (0,0)
	Person string
	X      int
	Y      int
}

type Child struct {
	Name      string
	Behaviour BehaviourType
	Address   Address
	WishList  []Present
	Presents  []Present
}
