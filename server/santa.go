package main

import (
	"flag"
	"fmt"
	"github.com/ChrisGora/semaphore"
	"math"
	"math/rand"
	"net"
	"net/rpc"
	"sync"
	"time"
	"workshop/util"
)

type SantaOperations struct {
	Workshop *rpc.Client
}

var numOfElves = 8
var semStorageRoom = semaphore.Init(4, 4)
var mTasks sync.Mutex
var mWorkshop sync.Mutex
var th util.TimeHandler

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
				time.Sleep(time.Duration(3 * len(child.WishList)) * time.Second)
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

// production - sends a request to the WorkshopOperations to begin working on the presents
func production(th *util.TimeHandler, inputChildren []util.Child, out chan []util.Child) {
	fmt.Println("Time:",th.GetTime(),
		fmt.Sprintf("Received work from Santa for %d children!", len(inputChildren)))

	numOfChildren := len(inputChildren) // needed as req.ChildrenList gets modified, so len(req.ChildrenList) won't work

	elves := make([]chan util.Child, numOfElves)
	for e := 0; e < numOfElves; e++ {
		elves[e] = make(chan util.Child)
	}

	for i := range elves {
		go elf(i, &inputChildren, elves[i])
	}

	outputChildren := []util.Child{}
	// Whichever elf returns some work, append to outputChildren, until all the presents have been made
	for len(outputChildren) < numOfChildren {
		select {
		case child := <-elves[0]:
			outputChildren = append(outputChildren, child)
		case child := <-elves[1]:
			outputChildren = append(outputChildren, child)
		case child := <-elves[2]:
			outputChildren = append(outputChildren, child)
		case child := <-elves[3]:
			outputChildren = append(outputChildren, child)
		case child := <-elves[4]:
			outputChildren = append(outputChildren, child)
		case child := <-elves[5]:
			outputChildren = append(outputChildren, child)
		case child := <-elves[6]:
			outputChildren = append(outputChildren, child)
		case child := <-elves[7]:
			outputChildren = append(outputChildren, child)
		}
	}

	fmt.Println("Time:",th.GetTime(),
		fmt.Sprintf("Completed work from Santa for %d children!", len(outputChildren)))

	out <- outputChildren
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
		if !util.ChildInSlice(child, *path) {
			xSquared := math.Pow(math.Abs(float64((*location).X - child.Address.X)), 2)
			ySquared := math.Pow(math.Abs(float64((*location).Y - child.Address.Y)), 2)
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
	path := []util.Child{}
	pathCalculated := make(chan bool)
	santasWorkshopLocation := util.Address{"Santa", 0, 0}
	//done := make(chan bool)
	//ticker := time.NewTicker(2 * time.Second)
	results := []util.Child{}
	children := make([]util.Child, len(req.ChildrenList))
	copy(children, req.ChildrenList)
	go production(&th, children, out)
	go calculatePath(req.ChildrenList, &path, &util.Address{"Santa", 0, 0}, pathCalculated)
	//go func(done chan bool, ticker *time.Ticker){
	//	for {
	//		select {
	//		case <-done:
	//			return
	//		case <-ticker.C:
	//			// TODO: make multiple messages to randomly print while he is waiting
	//			fmt.Println("Time:",th.GetTime(),"Santa drinks a cup of tea")
	//	}
	//}
	//}(done, ticker)

	<-pathCalculated
	// Note: path stores the children, route stores the addresses
	route := []util.Address{santasWorkshopLocation}
	for c := range path {
		route = append(route, path[c].Address)
	}
	route = append(route, santasWorkshopLocation)
	fmt.Println("Time:",th.GetTime(), "Santa's Route:", route)
	results = <-out
	//done <- true

	res.ChildrenList = results
	res.Route        = route

	fmt.Println("Time:",th.GetTime(),
		fmt.Sprintf("The presents for all %d children are ready to be delivered!", len(res.ChildrenList)))

	fmt.Println("##################################################")

	return
}

// main - registers RPC procedures
func main() {
	// Listen to client
	port := flag.String("client", "8080", "Port to listen on for client")
	flag.Parse()
	fmt.Println("Santa's Workshop is up and running!")

	rand.Seed(time.Now().UnixNano())
	santa := SantaOperations{}
	rpc.Register(&santa)

	client, err := net.Listen("tcp", ":"+*port)
	util.Check(err)

	defer client.Close()
	rpc.Accept(client)
}