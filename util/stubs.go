package util

var SantaHandler = "SantaOperations.Run"

type Response struct {
	ChildrenList  []Child
	Route         []Address
}

type Request struct {
	ChildrenList  []Child
}