package util

var WorkshopHandler = "WorkshopOperations.Run"

type Response struct {
	ChildrenList  []Child
	Route         []Address
}

type Request struct {
	Sender        int
	ChildrenList  []Child
}