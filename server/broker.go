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

var mWorkers sync.Mutex

var th util.TimeHandler

type WorkshopOperations struct {
	Workers []string
}

func delegateWork(worker *rpc.Client, work []util.Child, out chan []util.Child) {
	fmt.Println("Time:",th.GetTime(),"Sent message to an elf")
	request := util.Request{ChildrenList: work}
	response := new(util.Response)
	worker.Call(util.WorkerHandler, request, response)
	out <- response.ChildrenList
	worker.Close()
}

// Run - processes simulation of the workshop
func (workshop *WorkshopOperations) Run(req util.Request, res *util.Response) (err error) {
	th.SetStartTime()

	fmt.Println("Time:",th.GetTime(),
		fmt.Sprintf("Santa has received the wishlists of %d children!", len(req.ChildrenList)))

	channels := make([]chan []util.Child, len(workshop.Workers))
	dialedWorkers := []*rpc.Client{}
	w := 0
	for w < len(workshop.Workers) { // need to use a while loop with a counter because of indexing issues when workers disconnect
		worker, err := rpc.Dial("tcp", workshop.Workers[w])
		if err == nil {
			channels[w] = make(chan []util.Child)
			dialedWorkers = append(dialedWorkers, worker)
			w += 1
		} else {
			removeWorker(workshop, w)
		}
	}

	for i := 0; i < len(dialedWorkers); i++ {
		// Split work
		start := i * len(req.ChildrenList) / len(dialedWorkers)
		end := (i+1) * len(req.ChildrenList) / len(dialedWorkers)
		work := req.ChildrenList[start:end]
		go delegateWork(dialedWorkers[i], work, channels[i])
	}

	// TODO: after a certain amount of time, presume that one of the workers has stopped working
	// 	and therefore broker will need to send work to workers again
	res.ChildrenList = []util.Child{} // Remove all the children from before, when child.Presents was empty
	for i := 0; i < len(channels); i++ {
		// Note: depending on how you want to retrieve the results, you may need to rewrite this
		child := <-channels[i]
		res.ChildrenList = append(res.ChildrenList, child...)
	}

	fmt.Println("Time:",th.GetTime(),
		fmt.Sprintf("The presents for all %d children are ready to be delivered!", len(res.ChildrenList)))

	return
}

// addWorkers - adds/removes workers when a worker joins the server
func addWorkers(workshop *WorkshopOperations, listener net.Listener) {
	fmt.Println("Listening for workers")
	for {
		worker, err := listener.Accept()
		util.Check(err)
		wAddr := worker.RemoteAddr().String()

		// Send IP address back to the worker, so it can listen from that port for RPC calls from server
		wIP := strings.Split(wAddr, ":")[0]
		reader := bufio.NewReader(worker)
		wPort, err := reader.ReadString('\n')
		util.Check(err)
		worker.Close() // only needed to get the IP address

		// Difference between wAddr and workerAddr is that workerAddr is the address the worker listens on
		// while wAddr is the address that the worker uses to connect to the server
		workerAddr := wIP+":"+wPort[:len(wPort)-1] // need to remove \n at end of wPort
		fmt.Fprintln(worker, workerAddr)

		mWorkers.Lock()
		workshop.Workers = append(workshop.Workers, workerAddr)
		fmt.Println("Number of workers:", len(workshop.Workers))
		mWorkers.Unlock()
	}
}

// removeWorker - removes a worker from the list
func removeWorker(workshop *WorkshopOperations, w int) {
	mWorkers.Lock()
	workshop.Workers = append(workshop.Workers[:w], workshop.Workers[w+1:]...)
	mWorkers.Unlock()
	fmt.Println(fmt.Sprintf("Worker %d disconnected", w))
}

// main - registers RPC procedures
func main() {
	// Listen to client
	clientListenerPort := flag.String("client", "8080", "Port to listen on for client")
	workerListenerPort := flag.String("worker", "8081", "Port to listen on for workers")
	flag.Parse()
	fmt.Println("Broker running...")

	rand.Seed(time.Now().UnixNano())
	workshop := WorkshopOperations{}
	rpc.Register(&workshop)

	// Start goroutine to repeatedly accept workers joining server
	workerListener, err := net.Listen("tcp", ":"+*workerListenerPort)
	util.Check(err)
	defer workerListener.Close()
	go addWorkers(&workshop, workerListener)

	clientListener, err := net.Listen("tcp", ":"+*clientListenerPort)
	util.Check(err)

	defer clientListener.Close()
	rpc.Accept(clientListener)
}