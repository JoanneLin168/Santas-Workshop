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

const (
	screenWidth  = 640
	screenHeight = 480
)

var (
	santaImg    *ebiten.Image
	elfImg      *ebiten.Image
	workshopImg *ebiten.Image
)

// Game implements ebiten.Game interface.
type Game struct{}

// Update proceeds the game state.
// Update is called every tick (1/60 [s] by default).
func (g *Game) Update() error {
	// Write your game's logical update.
	// TODO: store the print statements into some queue here for visualisation
	return nil
}

// Draw draws the game screen.
// Draw is called every frame (typically 1/60[s] for 60Hz display).
func (g *Game) Draw(screen *ebiten.Image) {
	// Write your game's rendering.
	// TODO: move 8 elves to bottom (and store their orig pos), and move workshop to either top right or middle, whatever is easiest
	santaImgOp := &ebiten.DrawImageOptions{}
	santaImgOp.GeoM.Scale(0.2, 0.2)
	santaImgOp.GeoM.Translate(0, 0)
	screen.DrawImage(santaImg, santaImgOp)

	elfW, _ := elfImg.Size()
	elfImgOp := &ebiten.DrawImageOptions{}
	elfImgOp.GeoM.Scale(0.1, 0.1)
	elfImgOp.GeoM.Translate(float64(screenWidth-(elfW/10)), 0)
	screen.DrawImage(elfImg, elfImgOp)

	workshopW, workshopH := elfImg.Size()
	workshopImgOp := &ebiten.DrawImageOptions{}
	workshopImgOp.GeoM.Scale(0.2, 0.2)
	workshopImgOp.GeoM.Translate(float64((screenWidth/2)-((workshopW/5)/2)), float64(screenHeight-(workshopH/5)))
	screen.DrawImage(workshopImg, workshopImgOp)
}

// Layout takes the outside size (e.g., the window size) and returns the (logical) screen size.
// If you don't have to adjust the screen size with the outside size, just return a fixed size.
func (g *Game) Layout(outsideWidth, outsideHeight int) (sw, sh int) {
	return screenWidth, screenHeight
}

func read(conn *net.Conn, idChan chan int, done chan bool) {
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
		// TODO: will need to uncomment the ones below for when graphical window is implemented
		case "START":
			fmt.Fprintln(*conn, fmt.Sprintf("CLIENT-%d:OK", id))
		case "STOP":
			fmt.Fprintln(*conn, fmt.Sprintf("CLIENT-%d:OK", id))
		case "ELF_ENTER":
			fmt.Fprintln(*conn, fmt.Sprintf("CLIENT-%d:OK", id))
		case "ELF_EXIT":
			fmt.Fprintln(*conn, fmt.Sprintf("CLIENT-%d:OK", id))
		case "ROUTE":
			fmt.Fprintln(*conn, fmt.Sprintf("CLIENT-%d:OK", id))
		default: // TODO: have more messages sent back to server depending on situation, e.g. a message isn't received correctly idk
			fmt.Fprintln(*conn, fmt.Sprintf("CLIENT-%d:ERROR", id)) // TODO: update this and same in server
		}
	}
}

func decodeImages() {
	var err error
	santaImg, _, err = ebitenutil.NewImageFromFile("sprites/santa.png")
	util.Check(err)

	elfImg, _, err = ebitenutil.NewImageFromFile("sprites/elf.png")
	util.Check(err)

	workshopImg, _, err = ebitenutil.NewImageFromFile("sprites/workshop.png")
	util.Check(err)
}

func main() {
	// RPC dial to server
	server := flag.String("server", "127.0.0.1:8080", "IP:port string to connect to for RPC")
	client, err := rpc.Dial("tcp", *server)
	flag.Parse()
	util.Check(err)
	defer client.Close()

	// Set up IO
	idChan := make(chan int)
	done := make(chan bool)
	var conn net.Conn
	go func() {
		ioAddr := flag.String("ip", "127.0.0.1:8081", "IP:port string to connect to for bufio")
		flag.Parse()
		conn, err = net.Dial("tcp", *ioAddr)
		util.Check(err)
		read(&conn, idChan, done)
	}()
	id := <-idChan

	children := util.ConvertCSV("input.csv")

	decodeImages()

	go func() {
		results := []util.Child{}
		route := []util.Address{}
		time.Sleep(3 * time.Second)
		c.Run(id, client, children, &results, &route)
		fmt.Fprintln(conn, fmt.Sprintf("CLIENT-%d:CLOSE", id))
		<-done
	}()

	game := &Game{}
	// Specify the window size as you like. Here, a doubled size is specified.
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Santa's Workshop")
	// Call ebiten.RunGame to start your game loop.
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
