package util

var SantaHandler = "SantaOperations.Run"
var WorkshopHandler = "WorkshopOperations.Workshop"

type Response struct {
	ChildrenList  []Child
	Route         []Address
}

type Request struct {
	ChildrenList  []Child
}