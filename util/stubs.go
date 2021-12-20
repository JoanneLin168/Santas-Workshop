package util

var WorkshopHandler = "WorkshopOperations.Run"
var WorkerHandler = "WorkerOperations.Worker"

type Response struct {
	ChildrenList  []Child
}

type Request struct {
	ChildrenList  []Child
}