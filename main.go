package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	_ "image/png"
	"log"
	"net"
	"net/rpc"
	"strconv"
	"strings"
	"time"
	c "workshop/client"
	"workshop/util"
)

func success(conn *net.Conn, g *c.Game, task c.Task) {
	// close the channel to prevent blocking
	g.VisQueue <- task
	fmt.Fprintln(*conn, fmt.Sprintf("CLIENT-%d:OK", task.Id))
}

func read(conn *net.Conn, g *c.Game, idChan chan int, done chan bool) {
	reader := bufio.NewReader(*conn)
	id := -1
	for {
		msg, err := reader.ReadString('\n')
		util.Check(err)
		msg = msg[:len(msg)-1]
		msgSlice := strings.Split(msg, ":")
		command := msgSlice[0]
		message := msgSlice[1]
		fmt.Println(msg)

		switch command {
		case "ID":
			id, err = strconv.Atoi(message)
			util.Check(err)
			idChan <- id
			fmt.Fprintln(*conn, fmt.Sprintf("CLIENT-%d:OK", id))
		case "CLOSED":
			done <- true
			break
		case "START":
			task := c.Task{0, util.START, message}
			success(conn, g, task)
		case "STOP":
			task := c.Task{0, util.STOP, message}
			success(conn, g, task)
		case "ELF_ENTER":
			contentSlice := strings.Split(message, ";")
			id, err := strconv.Atoi(contentSlice[0])
			util.Check(err)
			child := contentSlice[1]
			task := c.Task{id, util.ELF_ENTER, child}
			success(conn, g, task)
		case "ELF_EXIT":
			contentSlice := strings.Split(message, ";")
			id, err := strconv.Atoi(contentSlice[0])
			util.Check(err)
			child := contentSlice[1]
			task := c.Task{id, util.ELF_EXIT, child}
			success(conn, g, task)
		case "ROUTE":
			task := c.Task{0, util.ROUTE, message} // TODO: consider displaying Santa's route visually
			success(conn, g, task)
		default: // TODO: have more messages sent back to server depending on situation, e.g. a message isn't received correctly idk
			fmt.Fprintln(*conn, fmt.Sprintf("CLIENT-%d:ERROR", id))
		}
	}
}

func main() {
	// RPC dial to server
	server := flag.String("server", "127.0.0.1:8080", "IP:port string to connect to for RPC")
	client, err := rpc.Dial("tcp", *server)
	flag.Parse()
	util.Check(err)
	defer client.Close()

	// Set up
	visQueue := make(chan c.Task, 100)
	workshopSpace := []int{-1, -1, -1, -1}
	start := make(chan bool)
	game := &c.Game { // Set up game
		Children:      []util.Child{},
		Completed:     map[string]bool{},
		VisSanta:      c.VisSprite{},
		VisElves:      []c.VisSprite{},
		Addresses:     []util.Address{},
		Route:         []util.Address{},
		VisQueue:      visQueue,
		WorkshopSpace: workshopSpace,
		Start:         start,
		Stage:         c.STANDBY,
	}
	ebiten.SetWindowSize(c.ScreenWidth, c.ScreenHeight)
	ebiten.SetWindowTitle("Santa's Workshop")

	// Set up IO
	idChan := make(chan int)
	done := make(chan bool)
	var conn net.Conn
	go func() {
		ioAddr := flag.String("ip", "127.0.0.1:8081", "IP:port string to connect to for bufio")
		flag.Parse()
		conn, err = net.Dial("tcp", *ioAddr)
		util.Check(err)
		read(&conn, game, idChan, done)
	}()
	id := <-idChan

	c.Init()

	// Initialise elves
	elfW, elfH := c.ElfImg.Size()
	for i := 0; i < 8; i++ {
		x := float64(2 * i * elfW)
		y := float64(c.ScreenHeight - 2*elfH)
		img, _, err := ebitenutil.NewImageFromFile("sprites/elf_"+strconv.Itoa(i)+".png")
		util.Check(err)
		c.MVisElves.Lock()
		elf := c.VisSprite{i, 0, c.STAND, x, y, x, y, x, y, img}
		game.VisElves = append(game.VisElves, elf)
		c.MVisElves.Unlock()
	}
	// Initialise santa
	x := float64(0)
	y := float64(0)
	game.VisSanta = c.VisSprite{0, 0, c.STAND, x, y, x, y, x, y, c.SantaImg}

	// RPC dial to server
	go func() {
		for {
			<-start // will start when Enter is pressed
			children := util.ConvertCSV("input.csv") // by adding it here, means you will use latest version of csv
			game.Children = children
			results := []util.Child{}
			route := []util.Address{}
			time.Sleep(2 * time.Second)
			c.Run(id, client, children, &results, &route)
		}
	}()

	// Run window
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	} else {
		fmt.Fprintln(conn, fmt.Sprintf("CLIENT-%d:CLOSE", id))
		<-done
	}
}
