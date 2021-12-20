package main

import (
	"flag"
	"fmt"
	"net"
	"net/rpc"
	"time"
	"workshop/util"
)

type WorkerOperations struct {}

// Worker - processes simulation of the workshop
func (w *WorkerOperations) Worker (req util.Request, res *util.Response) (err error) {
	var th util.TimeHandler
	th.SetStartTime()
	fmt.Println("Time:",th.GetTime(),
		fmt.Sprintf("Received work from Santa for %d children!", len(req.ChildrenList)))

	// TODO: calculate children list correctly - look at iPad for notes
	childrenList := []util.Child{}
	for c := range(req.ChildrenList) {
		child := req.ChildrenList[c]
		if req.ChildrenList[c].Behaviour == util.Good {
			time.Sleep(3 * time.Second)
			child.Presents = child.WishList
		} else {
			time.Sleep(1 * time.Second)
			child.Presents = append(child.Presents, util.Present{Type: util.Coal})
		}
		childrenList = append(childrenList, child)

	}
	res.ChildrenList = childrenList
	fmt.Println("Time:",th.GetTime(),"Completed work from Santa!")

	return
}

func main() {
	// Dial to server
	port := flag.String("port", "8030", "Port to listen on")
	serverPort := flag.String("server", "127.0.0.1:8081", "IP:port string to connect to")
	flag.Parse()
	server, err := net.Dial("tcp", *serverPort)
	util.Check(err)
	defer server.Close()
	fmt.Println("Connected to server")

	// Listens for RPC calls from server
	listener, err := net.Listen("tcp", ":"+*port)
	util.Check(err)
	defer listener.Close()

	// Sends IP:port to server to
	fmt.Fprintln(server, *port)
	worker := WorkerOperations{}
	rpc.Register(&worker)
	rpc.Accept(listener)
}
