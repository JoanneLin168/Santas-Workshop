package main

import (
	"flag"
	"fmt"
	"github.com/ChrisGora/semaphore"
	"net"
	"net/rpc"
	"sync"
	"time"
	"workshop/util"
)

type WorkshopOperations struct {}

var numOfElves = 8
var semStorageRoom = semaphore.Init(4, 4)
var mTasks sync.Mutex

func elf (id int, childrenList *[]util.Child, ch chan util.Child) {
	// Only one elf can access childrenList at a time to prevent any race conditions
	for {
		mTasks.Lock()
		if len(*childrenList) > 0 {
			child := (*childrenList)[0]
			*childrenList = (*childrenList)[1:]
			mTasks.Unlock()

			semStorageRoom.Wait()
			if child.Behaviour == util.Good {
				time.Sleep(3 * time.Second)
				child.Presents = child.WishList
			} else {
				time.Sleep(1 * time.Second)
				child.Presents = append(child.Presents, util.Present{Type: util.Coal})
			}

			ch <- child
			semStorageRoom.Post()
		} else {
			mTasks.Unlock()
			break
		}
	}

}

// Workshop - processes simulation of the workshop
func (workshop *WorkshopOperations) Workshop(req util.Request, res *util.Response) (err error) {
	var th util.TimeHandler
	th.SetStartTime()
	fmt.Println("Time:",th.GetTime(),
		fmt.Sprintf("Received work from Santa for %d children!", len(req.ChildrenList)))

	numOfChildren := len(req.ChildrenList) // needed as req.ChildrenList gets modified, so len(req.ChildrenList) won't work

	elves := make([]chan util.Child, numOfElves)
	for e := 0; e < numOfElves; e++ {
		elves[e] = make(chan util.Child)
	}

	for i := range(elves) {
		go elf(i, &req.ChildrenList, elves[i])
	}

	childrenList := []util.Child{}
	// Whichever elf returns some work, append to childrenList, until all the presents have been made
	for len(childrenList) < numOfChildren {
		select {
		case child := <-elves[0]:
			childrenList = append(childrenList, child)
		case child := <-elves[1]:
			childrenList = append(childrenList, child)
		case child := <-elves[2]:
			childrenList = append(childrenList, child)
		case child := <-elves[3]:
			childrenList = append(childrenList, child)
		case child := <-elves[4]:
			childrenList = append(childrenList, child)
		case child := <-elves[5]:
			childrenList = append(childrenList, child)
		case child := <-elves[6]:
			childrenList = append(childrenList, child)
		case child := <-elves[7]:
			childrenList = append(childrenList, child)
		}
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
