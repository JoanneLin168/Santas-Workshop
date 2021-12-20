package main

import (
	"flag"
	"fmt"
	"net"
	"net/rpc"
	"time"
	"workshop/util"
	"github.com/ChrisGora/semaphore"
)

type WorkshopOperations struct {}

var numOfElves = 8
var semElves = semaphore.Init(2, 2)

func elf (work []util.Child, ch chan []util.Child) {
	// work is the children list passed from the workshop, completed is the children list with presents to return
	semElves.Wait()
	completed := []util.Child{}
	for c := range(work) {
		child := work[c]
		if child.Behaviour == util.Good {
			time.Sleep(3 * time.Second)
			child.Presents = child.WishList
		} else {
			time.Sleep(1 * time.Second)
			child.Presents = append(child.Presents, util.Present{Type: util.Coal})
		}
		completed = append(completed, child)
	}
	ch <- completed
	semElves.Post()
}

// Workshop - processes simulation of the workshop
func (workshop *WorkshopOperations) Workshop(req util.Request, res *util.Response) (err error) {
	var th util.TimeHandler
	th.SetStartTime()
	fmt.Println("Time:",th.GetTime(),
		fmt.Sprintf("Received work from Santa for %d children!", len(req.ChildrenList)))

	elves := make([]chan []util.Child, numOfElves)
	for e := 0; e < numOfElves; e++ {
		elves[e] = make(chan []util.Child)
	}

	elvesWithWork := []int{}
	for i := 0; i < len(elves); i++ {
		// Split work
		start := i * len(req.ChildrenList) / len(elves)
		end := (i+1) * len(req.ChildrenList) / len(elves)
		work := req.ChildrenList[start:end]

		if len(work) > 0 { // prevents elves with no work from doing work
			go elf(work, elves[i])
			elvesWithWork = append(elvesWithWork, i)
		}
	}

	childrenList := []util.Child{}
	for i := range(elvesWithWork) {
		index := elvesWithWork[i]
		childrenList = append(childrenList, <-elves[index]...)
	}

	res.ChildrenList = childrenList
	fmt.Println("Time:",th.GetTime(),
		fmt.Sprintf("Completed work from Santa for %d children!", len(res.ChildrenList)))

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
	workshop := WorkshopOperations{}
	rpc.Register(&workshop)
	rpc.Accept(listener)
}
