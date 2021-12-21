package main

import (
	"bufio"
	"flag"
	"fmt"
	"math/rand"
	"net"
	"net/rpc"
	"strings"
	"sync"
	"time"
	"workshop/util"
)

var mWorkshop sync.Mutex

var th util.TimeHandler

type SantaOperations struct {
	Workshop *rpc.Client
}

// beginProduction - sends a request to the WorkshopOperations to begin working on the presents
func beginProduction(workshop *rpc.Client, children []util.Child, out chan []util.Child) {
	request := util.Request{ChildrenList: children}
	response := new(util.Response)
	workshop.Call(util.WorkshopHandler, request, response)
	out <- response.ChildrenList
}

// Run - processes simulation of the workshop
func (santa *SantaOperations) Run(req util.Request, res *util.Response) (err error) {
	th.SetStartTime()

	fmt.Println("Time:",th.GetTime(),
		fmt.Sprintf("Santa has received the wishlists of %d children!", len(req.ChildrenList)))

	// if there are no children in the request, just return an empty response
	if len(req.ChildrenList) == 0 {
		return
	}

	out := make(chan []util.Child)
	ticker := time.NewTicker(2 * time.Second)
	done := make(chan bool)
	results := []util.Child{}
	go beginProduction(santa.Workshop, req.ChildrenList, out)
	go func(done chan bool, ticker *time.Ticker){
		for {
			select {
			case <-done:
				return
			case <-ticker.C:
				// TODO: make multiple messages to randomly print while he is waiting
				fmt.Println("Time:",th.GetTime(),"Santa drinks a cup of tea")
		}
	}
	}(done, ticker)
	results = <-out
	done <- true

	// TODO: after a certain amount of time, presume that one of the workers has stopped working
	// 		and therefore broker will need to send work to workers again
	//		Note: not too sure if this is still applicable
	res.ChildrenList = results

	fmt.Println("Time:",th.GetTime(),
		fmt.Sprintf("The presents for all %d children are ready to be delivered!", len(res.ChildrenList)))

	// TODO: Santa figures out a path to deliver to all of the children - perhaps greedy algorithm?
	// 	Note: this is the travelling salesman problem

	return
}

func connectToWorkshop(santa *SantaOperations, listener net.Listener) {
	fmt.Println("Waiting for Workshop")

	for {
		workshopReceiver, err := listener.Accept()
		util.Check(err)

		wAddr := workshopReceiver.RemoteAddr().String()

		// Send IP address back to the worker, so it can listen from that port for RPC calls from server
		wIP := strings.Split(wAddr, ":")[0]
		reader := bufio.NewReader(workshopReceiver)
		wPort, err := reader.ReadString('\n')
		util.Check(err)
		workshopReceiver.Close() // only needed to get the IP address

		// Difference between wAddr and workerAddr is that workerAddr is the address the worker listens on
		// while wAddr is the address that the worker uses to connect to the server
		workshopAddr := wIP+":"+wPort[:len(wPort)-1] // need to remove \n at end of wPort
		fmt.Fprintln(workshopReceiver, workshopAddr)

		mWorkshop.Lock()
		workshopSender, err := rpc.Dial("tcp", workshopAddr)
		// TODO: make sure that this has no errors, e.g. while Santa is waiting on the workshop to finish making the presents,
		// 		and a new workshop is added, make sure to resend the work to the new workshop (means it will work both if the
		//		previous workshop is still alive AND if it has crashed, and a new one started)
		santa.Workshop = workshopSender // Note: this will replace the previous workshop
		fmt.Println("Connected to workshop")
		mWorkshop.Unlock()
	}
}

// main - registers RPC procedures
func main() {
	// Listen to client
	clientListenerPort := flag.String("client", "8080", "Port to listen on for client")
	workerListenerPort := flag.String("worker", "8081", "Port to listen on for workers")
	flag.Parse()
	fmt.Println("Broker running...")

	rand.Seed(time.Now().UnixNano())
	santa := SantaOperations{}
	rpc.Register(&santa)

	// Start goroutine to repeatedly accept workers joining server
	workerListener, err := net.Listen("tcp", ":"+*workerListenerPort)
	util.Check(err)
	defer workerListener.Close()
	go connectToWorkshop(&santa, workerListener)

	clientListener, err := net.Listen("tcp", ":"+*clientListenerPort)
	util.Check(err)

	defer clientListener.Close()
	rpc.Accept(clientListener)
}