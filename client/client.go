package client

import (
	"fmt"
	"net/rpc"
	"workshop/util"
)

func Run(id int, client *rpc.Client, children []util.Child, results *[]util.Child, route *[]util.Address) {
	request := util.Request{id, children}
	response := new(util.Response)
	var th util.TimeHandler
	th.SetStartTime()
	fmt.Println("Start Time:",th.GetTime(),
		fmt.Sprintf("Sent Santa the wishlists of %d children!", len(children)))

	client.Call(util.WorkshopHandler, request, response)

	fmt.Println("Duration time:",th.GetTime(),
		fmt.Sprintf("All %d children has received presents from Santa!", len(children)))

	*results = response.ChildrenList
	*route   = response.Route
}
