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
	//Child Child <- not too sure if needed or not
	Type  presentType
}

type behaviourType uint8
const (
	Good behaviourType = iota
	Bad
)
type Child struct {
	Name  string
	Behaviour behaviourType
	WishList []Present
	Presents []Present
}
