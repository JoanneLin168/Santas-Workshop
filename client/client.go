package client

import (
	"fmt"
	"net/rpc"
	"workshop/util"
)

func Run(client *rpc.Client, children []util.Child, results *[]util.Child) {
	request := util.Request{ChildrenList: children}
	response := new(util.Response)
	var th util.TimeHandler
	th.SetStartTime()
	fmt.Println("Start Time:",th.GetTime(),"Sent Santa the wishlists!")
	client.Call(util.SantaHandler, request, response)

	// TODO: create unit tests for testing with a list of children and their presents
	// Debugging purposes, makes sure the gifts are correct
	fmt.Println("Duration:",th.GetTime(),"Every child has received gifts from Santa!")

	*results = response.ChildrenList
}
