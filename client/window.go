package client

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/examples/resources/fonts"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
	"image/color"
	"log"
	"strconv"
	"strings"
	"sync"
	"workshop/util"
)

var MVisElves sync.Mutex
var mWorkshop sync.Mutex

type Task struct {
	Id      int // mainly used for elves
	Action  util.ActionType
	Content string
}

var (
	SantaImg    *ebiten.Image
	ElfImg      *ebiten.Image
	WorkshopImg *ebiten.Image
	Font        font.Face

	length  = float64(128) // Note: length of one quadrant, not the entire axis
	OriginX = float64(ScreenWidth/2)
	OriginY = float64(144) + length
)

const (
	ScreenWidth  = 640
	ScreenHeight = 480
	FontSize     = 12
	MapUnit      = 128/5
	MapPoint     = 8
)

type MovementType uint8
const (
	STAND MovementType = iota
	ENTER
	EXIT
)

type Stages uint8
const (
	STANDBY Stages = iota
	PROCESSING
	COMPLETED
)

// VisSprite - stores the id and the position of an elf sprite
type VisSprite struct {
	Id     int
	Frame  int
	Move   MovementType
	X      float64
	Y      float64
	StartX float64
	StartY float64
	EndX   float64
	EndY   float64
	Img    *ebiten.Image
}

// drawElves - draws elves
func (s VisSprite) drawElves(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(2, 2)

	op.GeoM.Translate(s.X, s.Y)
	screen.DrawImage(s.Img, op)
}

// drawSanta - draws Santa (has different logic)
func (s VisSprite) drawSanta(screen *ebiten.Image) {
	w, h := SantaImg.Size()
	op := &ebiten.DrawImageOptions{}
	x := OriginX + (s.X * MapUnit) - float64(w/2)
	y := OriginY + (s.Y * MapUnit) - float64(h/2)

	op.GeoM.Translate(x, y)
	screen.DrawImage(s.Img, op)
}

// Game implements ebiten.Game interface.
type Game struct {
	Children      []util.Child
	Completed     map[string]bool
	VisSanta      VisSprite
	VisElves      []VisSprite
	Addresses     []util.Address
	Route         []util.Address
	VisQueue      chan Task
	Log           []string
	WorkshopSpace []int
	Start         chan bool
	Stage         Stages
}

// updateElvesPos - updates the position of the elves every time Update() is called
func (g *Game) updateElvesPos() {
	MVisElves.Lock()
	for i := range g.VisElves {
		if g.VisElves[i].Frame < 30 && g.VisElves[i].Move == ENTER {
			g.VisElves[i].X += (g.VisElves[i].EndX - g.VisElves[i].StartX) / 30
			g.VisElves[i].Y += (g.VisElves[i].EndY - g.VisElves[i].StartY) / 30
			g.VisElves[i].Frame += 1
		} else if g.VisElves[i].Frame < 30 && g.VisElves[i].Move == EXIT {
			g.VisElves[i].X += (g.VisElves[i].StartX - g.VisElves[i].EndX) / 30
			g.VisElves[i].Y += (g.VisElves[i].StartY - g.VisElves[i].EndY) / 30
			g.VisElves[i].Frame += 1
		} else if g.VisElves[i].Move == ENTER {
			g.VisElves[i].X = g.VisElves[i].EndX
			g.VisElves[i].Y = g.VisElves[i].EndY
			g.VisElves[i].Frame = 0
			g.VisElves[i].Move = STAND
		} else if g.VisElves[i].Move == EXIT {
			g.VisElves[i].X = g.VisElves[i].StartX
			g.VisElves[i].Y = g.VisElves[i].StartY
			g.VisElves[i].Frame = 0
			g.VisElves[i].Move = STAND
		}
	}
	MVisElves.Unlock()
}

// updateSantaPos - updates the position of Santa every time Update() is called
func (g *Game) updateSantaPos() {
	// For santa, use Id to track which address he is moving to
	// Just use ENTER to say that he needs to move
	if g.VisSanta.Frame < 30 && g.VisSanta.Move == ENTER {
		g.VisSanta.X += (g.VisSanta.EndX - g.VisSanta.StartX) / 30
		g.VisSanta.Y += (g.VisSanta.EndY - g.VisSanta.StartY) / 30
		g.VisSanta.Frame += 1
	} else if g.VisSanta.Move == ENTER { // if frame == 30, update his next destination
		g.VisSanta.X = g.VisSanta.EndX
		g.VisSanta.Y = g.VisSanta.EndY
		g.VisSanta.StartX = g.VisSanta.X
		g.VisSanta.StartY = g.VisSanta.Y
		if g.VisSanta.Id + 1 < len(g.Route) {
			g.VisSanta.Id += 1
		} else {
			g.VisSanta.Id = 0
		}
		g.VisSanta.Frame = 0
		i := g.VisSanta.Id

		g.VisSanta.EndX = float64(g.Route[i].X)
		g.VisSanta.EndY = float64(g.Route[i].Y)
	}
}

// addTaskToLog - adds a task to the log in string form
func (g *Game) addTaskToLog(task Task) {
	str := ""
	switch task.Action {
	case util.START:
		str = "Santa's workshop has started working on the presents"
	case util.STOP:
		str = "Santa's workshop has completed the presents"
	case util.ELF_ENTER:
		str = fmt.Sprintf("Elf %d has entered the storage room to work on %s's presents", task.Id, task.Content)
	case util.ELF_EXIT:
		str = fmt.Sprintf("Elf %d has finished %s's presents", task.Id, task.Content)
	case util.ROUTE:
		str = "Santa has figured out his route for Christmas Eve"
	}
	if len(g.Log) < 10 {
		g.Log = append(g.Log, str)
	} else {
		g.Log = g.Log[1:]
		g.Log = append(g.Log, str)
	}
}

// drawMap - draws out a map with the positions of the children
func (g *Game) drawMap(screen *ebiten.Image) {
	// draw out the outline of the map
	ebitenutil.DrawLine(screen, OriginX-length, OriginY, OriginX+length, OriginY, color.White) // x-axis
	ebitenutil.DrawLine(screen, OriginX, OriginY-length, OriginX, OriginY+length, color.White) // y-axis

	// draw out the positions of the children
	if g.Stage != STANDBY && len(g.Addresses) < 1 {
		g.Addresses = append(g.Addresses, util.SantaAddr)
		for c := range g.Children {
			child := g.Children[c]
			g.Addresses = append(g.Addresses, child.Address)
		}
	}
	for a := range g.Addresses {
		addr := g.Addresses[a]
		x := float64(addr.X * MapUnit)
		y := float64(addr.Y * MapUnit)
		ebitenutil.DrawRect(screen, OriginX+x-(MapPoint/2), OriginY+y-(MapPoint/2), MapPoint, MapPoint, color.White)
		lengthOfName := text.BoundString(Font, addr.Person)
		text.Draw(screen, addr.Person, Font,
			int(OriginX+x)-(lengthOfName.Size().X/2), int(OriginY+y)-(lengthOfName.Size().Y/2), color.White)
	}
}

func (g *Game) drawRoute(screen *ebiten.Image) {
	for a := 0; a < len(g.Route)-1; a++ {
		curr := g.Route[a]
		next := g.Route[a+1]
		x1 := float64(curr.X * MapUnit)
		y1 := float64(curr.Y * MapUnit)
		x2 := float64(next.X * MapUnit)
		y2 := float64(next.Y * MapUnit)
		ebitenutil.DrawLine(screen, OriginX+x1, OriginY+y1, OriginX+x2, OriginY+y2, color.RGBA{255, 0, 0, 255})
	}
}

// Update proceeds the game state.
// Update is called every tick (1/60 [s] by default).
func (g *Game) Update() error {
	// Write your game's logical update.

	if g.Stage == STANDBY || g.Stage == COMPLETED {
		if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
			g.Start <- true
		}
	}

	g.updateElvesPos() // update positions of elves
	g.updateSantaPos() // update position of Santa
	select {
	case task := <-g.VisQueue: // upon receiving a task, update visualisation
		// Add task to Log
		g.addTaskToLog(task)

		// Handle task logic
		switch task.Action {
		case util.START:
			g.Stage = PROCESSING
			g.Completed = map[string]bool{}
			g.Log = []string{}
			g.Route = []util.Address{}
			g.Addresses = []util.Address{}

			// reset Santa
			x := float64(0)
			y := float64(0)
			g.VisSanta.Move = STAND
			g.VisSanta.X = x
			g.VisSanta.Y = y
			g.VisSanta.StartX = x
			g.VisSanta.StartY = y
			g.VisSanta.EndX = x
			g.VisSanta.EndY = y
			g.VisSanta.Frame = 0
		case util.STOP:
			g.Stage = COMPLETED
		case util.ELF_ENTER:
			workshopW, workshopH := WorkshopImg.Size()
			mWorkshop.Lock()
			for s := range g.WorkshopSpace {
				if g.WorkshopSpace[s] == -1 {
					g.WorkshopSpace[s] = task.Id
					MVisElves.Lock()
					if (s+1) % 2 == 0 { // if even, go to right
						g.VisElves[task.Id].EndX = float64(ScreenWidth-workshopW) // half of workshopW size (don't forget it is 2x scaled)
					} else {
						g.VisElves[task.Id].EndX = float64(ScreenWidth-2*workshopW)
					}
					if s+1 > 2 { // if > 2, go to bottom
						g.VisElves[task.Id].EndY = float64(workshopH) // half of workshopH size (don't forget it is 2x scaled)
					} else {
						g.VisElves[task.Id].EndY = float64(0)
					}
					g.VisElves[task.Id].Frame = 0
					g.VisElves[task.Id].Move = ENTER
					MVisElves.Unlock()
					break
				}
			}
			mWorkshop.Unlock()
		case util.ELF_EXIT:
			mWorkshop.Lock()
			for s := range g.WorkshopSpace {
				if g.WorkshopSpace[s] == task.Id {
					g.WorkshopSpace[s] = -1
					break
				}
			}
			mWorkshop.Unlock()
			MVisElves.Lock()
			g.VisElves[task.Id].Frame = 0
			g.VisElves[task.Id].Move = EXIT
			g.Completed[task.Content] = true
			MVisElves.Unlock()
		case util.ROUTE:
			str := task.Content[0:len(task.Content)-2]
			strSlice := strings.Split(str, "}")
			strSlice[len(strSlice)-1] += " " // adds space to be removed in for loop
			for s := range strSlice {
				strSlice[s] = strSlice[s][2:]
				details := strings.Split(strSlice[s], " ")
				x, err := strconv.Atoi(details[1])
				util.Check(err)
				y, err := strconv.Atoi(details[2])
				util.Check(err)
				addr := util.Address{Person: details[0], X: x, Y: y}
				g.Route = append(g.Route, addr)
			}
			g.VisSanta.Move = ENTER
			g.VisSanta.Id += 1
			g.VisSanta.EndX = float64(g.Route[g.VisSanta.Id].X)
			g.VisSanta.EndY = float64(g.Route[g.VisSanta.Id].Y)

		}
	default:
		// do nothing
	}

	return nil
}

// Draw draws the game screen.
// Draw is called every frame (typically 1/60[s] for 60Hz display).
func (g *Game) Draw(screen *ebiten.Image) {
	// Write your game's rendering.
	workshopW, workshopH := WorkshopImg.Size()
	workshopImgOp := &ebiten.DrawImageOptions{}
	workshopImgOp.GeoM.Scale(2, 2)
	workshopImgOp.GeoM.Translate(float64(ScreenWidth-2*workshopW), 0)
	screen.DrawImage(WorkshopImg, workshopImgOp)
	text.Draw(screen, "STORAGE ROOM", Font, ScreenWidth-(2*workshopW),
		2*workshopH+FontSize, color.White)

	for e := range g.VisElves {
		MVisElves.Lock()
		g.VisElves[e].drawElves(screen)
		MVisElves.Unlock()
	}

	// Text in the bottom-right to tell you what to do
	switch g.Stage {
	case STANDBY:
		msg := "Press 'Enter' to start"
		text.Draw(screen, msg, Font, ScreenWidth-128, ScreenHeight-FontSize, color.White)
	case PROCESSING:
		msg := "Processing..."
		text.Draw(screen, msg, Font, ScreenWidth-128, ScreenHeight-FontSize, color.White)
	case COMPLETED:
		msg := "Completed!"
		text.Draw(screen, msg, Font, ScreenWidth-128, ScreenHeight-FontSize, color.White)
		msg2 := "Press 'Enter' to start"
		text.Draw(screen, msg2, Font, ScreenWidth-128, ScreenHeight, color.White)
	}

	// Draw out map and route:
	if len(g.Route) > 0 {
		g.drawRoute(screen)
	}
	g.drawMap(screen)
	text.Draw(screen, "<- Santa's Route", Font, ScreenWidth-(2*workshopW)-56,
		int(OriginY)+(FontSize/3), color.White)

	// Draw Santa on top of map
	g.VisSanta.drawSanta(screen)

	// Draw out list of children
	listX := 0
	listY := int(OriginY - 2*length/3)
	text.Draw(screen, "LIST OF CHILDREN:", Font, listX, listY, color.White)
	if g.Stage != STANDBY {
		for c := range g.Children {
			name := (g.Children)[c].Name
			msg := name+"\n"
			y := listY + (c+1)*FontSize
			var msgColor color.Color
			if _, ok := (g.Completed)[name]; ok {
				msgColor = color.RGBA{0, 255, 0, 255}
			} else {
				msgColor = color.White
			}
			text.Draw(screen, msg, Font, listX, y, msgColor)
		}
	}

	// Log of what is happening
	text.Draw(screen, "LOG:", Font, 0, FontSize, color.White)
	for i := range g.Log {
		msg := g.Log[i]
		text.Draw(screen, msg, Font, 0, FontSize+(i+1)*FontSize, color.White)
	}
}

// Layout takes the outside size (e.g., the window size) and returns the (logical) screen size.
// If you don't have to adjust the screen size with the outside size, just return a fixed size.
func (g *Game) Layout(outsideWidth, outsideHeight int) (sw, sh int) {
	return ScreenWidth, ScreenHeight
}

// Init - prepares assets
func Init() {
	// Images
	var err error
	SantaImg, _, err = ebitenutil.NewImageFromFile("sprites/santa.png")
	util.Check(err)
	ElfImg, _, err = ebitenutil.NewImageFromFile("sprites/elf_0.png") // used only for the size
	util.Check(err)
	WorkshopImg, _, err = ebitenutil.NewImageFromFile("sprites/workshop.png")
	util.Check(err)

	// Font
	tt, err := opentype.Parse(fonts.MPlus1pRegular_ttf)
	if err != nil {
		log.Fatal(err)
	}
	const dpi = 72
	Font, err = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    FontSize,
		DPI:     dpi,
		Hinting: font.HintingFull,
	})
}