package main

import (
	"bufio"
	"flag"
	"fmt"
	"math"
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

// childInSlice - checks if a child is in the slice. Note: compares names - might need to make a comparison function for Child
func childInSlice(child util.Child, children []util.Child) bool {
	for c := range(children) {
		if child.Name == children[c].Name {
			return true
		}
	}
	return false
}

// calculatePath - Generate a path for Santa's route on Christmas Eve
func calculatePath(children []util.Child, path *[]util.Child, location *util.Address, pathCalculated chan bool) {
	// Santa is lazy, and so he wants to take the shortest path between every child - greedily!
	if len(children) == 0 {
		// Note: doesn't add Santa's Workshop as the final destination because of path's type.
		pathCalculated <- true
		return
	}

	distancesMap := make(map[float64]util.Child)
	for c := range children {
		child := (children)[c]
		if !childInSlice(child, *path) {
			xSquared := math.Abs(float64((*location).X ^ 2 - child.Address.Y ^ 2))
			ySquared := math.Abs(float64((*location).Y^2 - child.Address.Y^2))
			distance := math.Sqrt(xSquared+ySquared)

			// if distance was not a key in distancesMap, add it
			if _, ok := distancesMap[distance]; !ok {
				distancesMap[distance] = child
			}
		}
	}

	distancesKeys := []float64{}
	for k := range distancesMap {
		distancesKeys = append(distancesKeys, k)
	}
	shortestDistance := util.MinFloat64(distancesKeys)
	closestChild := distancesMap[shortestDistance]
	*path = append(*path, closestChild)

	index := -1
	for i := range children {
		if children[i].Name == closestChild.Name {
			index = i
			break
		}
	}

	*location = closestChild.Address
	remaining := util.RemoveChild(index, children)

	calculatePath(remaining, path, location, pathCalculated)
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
	done := make(chan bool)
	path := []util.Child{}
	pathCalculated := make(chan bool)
	santasWorkshopLocation := util.Address{"Santa", 0, 0}
	ticker := time.NewTicker(2 * time.Second)
	results := []util.Child{}
	go beginProduction(santa.Workshop, req.ChildrenList, out)
	go calculatePath(req.ChildrenList, &path, &util.Address{"Santa", 0, 0}, pathCalculated)
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

	<-pathCalculated
	// Note: path stores the children, route stores the addresses
	route := []util.Address{santasWorkshopLocation}
	for c := range path {
		route = append(route, path[c].Address)
	}
	route = append(route, santasWorkshopLocation)
	fmt.Println("Time:",th.GetTime(), "Santa's Route:", route)
	results = <-out
	done <- true

	// TODO: after a certain amount of time, presume that one of the workers has stopped working
	// 		and therefore broker will need to send work to workers again
	//		Note: not too sure if this is still applicable
	res.ChildrenList = results

	fmt.Println("Time:",th.GetTime(),
		fmt.Sprintf("The presents for all %d children are ready to be delivered!", len(res.ChildrenList)))

	fmt.Println("##################################################")

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