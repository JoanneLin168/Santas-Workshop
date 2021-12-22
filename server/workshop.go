package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/ChrisGora/semaphore"
	"math"
	"net"
	"net/rpc"
	"strconv"
	"strings"
	"sync"
	"time"
	"workshop/util"
)

type WorkshopOperations struct {
	Clients *map[int]net.Conn
}

type actionType uint8
const (
	START actionType = iota
	STOP
	ELF_ENTER
	ELF_EXIT
	ROUTE
)

type sendToClient struct {
	val bool
	client net.Conn
}

var numOfElves = 8
var semStorageRoom = semaphore.Init(4, 4)
var mTasks sync.Mutex

// log - logs what is happening in the system
func log(th *util.TimeHandler, toSend *sendToClient, str string, aType actionType) {
	/**
	START messages:
		- number of children in list
	STOP messages:
		- number of children in list
	ELF messages (include ID in every message):
		- elf enter storage room
		- elf exit storage room
	ROUTE messages
		- the path
	 */
	var action string
	switch aType {
	case START :
		action = "START"
	case STOP:
		action = "STOP"
	case ELF_ENTER:
		action = "ELF_ENTER"
	case ELF_EXIT:
		action = "ELF_EXIT"
	case ROUTE:
		action = "ROUTE"
	}
	toPrint := action+":"+str
	fmt.Println("Time:",th.GetTime(),toPrint)

	if toSend.val {
		fmt.Fprintln(toSend.client, toPrint)
	}

}

// elf - enters the storage room and prepares the present for each child
func elf (id int, toSend *sendToClient, th *util.TimeHandler, childrenList *[]util.Child, ch chan util.Child) {
	// Only one elf can access childrenList at a time to prevent any race conditions
	for {
		mTasks.Lock()
		if len(*childrenList) > 0 {
			child := (*childrenList)[0]
			*childrenList = (*childrenList)[1:]
			mTasks.Unlock()

			semStorageRoom.Wait()
			str := fmt.Sprintf("%d;%s", id, child.Name)
			log(th, toSend, str, ELF_ENTER)

			if child.Behaviour == util.Good {
				time.Sleep(time.Duration(3 * len(child.WishList)) * time.Second)
				child.Presents = child.WishList
			} else {
				time.Sleep(1 * time.Second)
				child.Presents = append(child.Presents, util.Present{Type: util.Coal})
			}
			ch <- child
			semStorageRoom.Post()
			str = fmt.Sprintf("%d;%s", id, child.Name)
			log(th, toSend, str, ELF_EXIT)
		} else {
			mTasks.Unlock()
			break
		}
	}

}

// production - sends a request to the WorkshopOperations to begin working on the presents
func production(toSend *sendToClient, th *util.TimeHandler, inputChildren []util.Child, out chan []util.Child) {
	numOfChildren := len(inputChildren) // needed as req.ChildrenList gets modified, so len(req.ChildrenList) won't work

	elves := make([]chan util.Child, numOfElves)
	for e := 0; e < numOfElves; e++ {
		elves[e] = make(chan util.Child)
	}

	for i := range elves {
		go elf(i, toSend, th, &inputChildren, elves[i])
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

	out <- outputChildren
}

// calculatePath - generate a path for Santa's route on Christmas Eve
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
func (santa *WorkshopOperations) Run(req util.Request, res *util.Response) (err error) {
	// TODO: make this wait until client is in the clients map, maybe?
	toSend := sendToClient{val: req.Sender != 0}
	if toSend.val { toSend.client = (*santa.Clients)[req.Sender] }
	var th util.TimeHandler
	th.SetStartTime()
	numOfChildren := len(req.ChildrenList)
	log(&th, &toSend, strconv.Itoa(numOfChildren), START)

	// if there are no children in the request, just return an empty response
	if numOfChildren == 0 {
		log(&th, &toSend, "0", STOP)
		return
	}

	out := make(chan []util.Child)
	path := []util.Child{}
	pathCalculated := make(chan bool)
	santasWorkshopLocation := util.Address{"Santa", 0, 0}
	results := []util.Child{}
	children := make([]util.Child, len(req.ChildrenList))
	copy(children, req.ChildrenList)

	go production(&toSend, &th, children, out)
	go calculatePath(req.ChildrenList, &path, &util.Address{"Santa", 0, 0}, pathCalculated)

	<-pathCalculated
	route := []util.Address{santasWorkshopLocation} // Note: path stores the children, route stores the addresses
	for c := range path {
		route = append(route, path[c].Address)
	}
	route = append(route, santasWorkshopLocation)
	str := fmt.Sprintf("%v", route)
	log(&th, &toSend, str, ROUTE)
	results = <-out

	res.ChildrenList = results
	res.Route = route
	log(&th, &toSend, strconv.Itoa(numOfChildren), STOP)

	fmt.Println("##################################################")

	return
}

// handleClient - when a client sends a message back, check that it is an ok message, and deal with it accordingly
func handleClient(client net.Conn, clientid int, msgs chan util.Message) {
	reader := bufio.NewReader(client)
	for {
		msg, err := reader.ReadString('\n')
		util.Check(err)
		msg = msg[:len(msg)-1]
		msgSlice := strings.Split(msg, ":")
		message := msgSlice[1]
		msgs <- util.Message{Sender: clientid, Message: message}
		if message == "CLOSE" {
			break
		}
	}
}

// acceptConns - continuously accept a network connection from the Listener and add it to the channel for handling connections.
func acceptConns(ln net.Listener, conns chan net.Conn) {
	for {
		conn, err := ln.Accept()
		util.Check(err)
		conns <- conn
	}
}

// main - handles listeners
func main() {
	rpcPort := flag.String("rpc", "8080", "Port to listen on for listener RPC")
	ioPort := flag.String("bufio", "8081", "Port to listen on for listener bufio")
	flag.Parse()
	fmt.Println("Santa's Workshop is up and running!")

	conns := make(chan net.Conn)
	msgs := make(chan util.Message)
	clients := make(map[int]net.Conn)

	// IO code
	go func(){
		ioListener, err := net.Listen("tcp", ":"+*ioPort)
		util.Check(err)
		go acceptConns(ioListener, conns)
		for {
			// TODO: handle clients disconnecting
			select {
			case conn := <-conns:
				id := len(clients) + 1 // id 0 is for the Workshop
				clients[id] = conn
				fmt.Println("Client", id, "has connected!")
				fmt.Fprintln(conn, "ID:"+strconv.Itoa(id))
				go handleClient(conn, id, msgs)
			case msg := <-msgs:
				// TODO: handle clients that don't send back an "ok" message after receiving from server
				// 	also handle when client doesn't send an "ok" message back after 10 seconds (tbc)
				if msg.Message == "OK" {
					continue
				} else if msg.Message == "CLOSE" {
					fmt.Fprintln(clients[msg.Sender], "CLOSED:0")
					delete(clients, msg.Sender)
				} else {
					fmt.Println("Uh oh...")
				}
			}
		}
	}()

	// RPC code
	santa := WorkshopOperations{&clients}
	rpc.Register(&santa)
	rpcListener, err := net.Listen("tcp", ":"+*rpcPort)
	util.Check(err)
	defer rpcListener.Close()
	rpc.Accept(rpcListener)
}