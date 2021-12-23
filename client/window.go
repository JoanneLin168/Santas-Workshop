package client

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/examples/resources/fonts"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
	"image/color"
	"log"
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
)

const (
	ScreenWidth  = 640
	ScreenHeight = 480
)

type ElfMovement uint8
const (
	STAND ElfMovement = iota
	ENTER
	EXIT
)

type Stages uint8
const (
	STANDBY Stages = iota
	PROCESSING
	COMPLETED
)

// VisElf - stores the id and the position of an elf sprite
type VisElf struct {
	Id     int
	Frame  int
	Move   ElfMovement
	X      float64
	Y      float64
	StartX float64
	StartY float64
	EndX   float64
	EndY   float64
	Img    *ebiten.Image
}
// drawElf - draws elves
func (e VisElf) drawElf(screen *ebiten.Image) {
	elfImgOp := &ebiten.DrawImageOptions{}
	elfImgOp.GeoM.Scale(2, 2)

	elfImgOp.GeoM.Translate(e.X, e.Y)
	screen.DrawImage(e.Img, elfImgOp)
}

// Game implements ebiten.Game interface.
type Game struct {
	VisElves      []VisElf
	VisRoute      string
	VisQueue      chan Task
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

	select {
	case task := <-g.VisQueue: // upon receiving a task, update visualisation
		switch task.Action {
		case util.START:
			g.Stage = PROCESSING
		case util.STOP:
			close(g.VisQueue)
			g.VisQueue = make(chan Task) // make a new channel so the old channel doesn't send zeroes
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
			MVisElves.Unlock()
		case util.ROUTE:
			str := task.Content[1:len(task.Content)-1]
			strSlice := strings.Split(str, "}")
			route := strings.Join(strSlice, "}\n")
			g.VisRoute = route
		}
	default: // do nothing, stops the function from blocking
	}

	return nil
}

// Draw draws the game screen.
// Draw is called every frame (typically 1/60[s] for 60Hz display).
func (g *Game) Draw(screen *ebiten.Image) {
	// Write your game's rendering.
	santaImgOp := &ebiten.DrawImageOptions{}
	santaImgOp.GeoM.Scale(2, 2)
	santaImgOp.GeoM.Translate(0, 0)
	screen.DrawImage(SantaImg, santaImgOp)

	switch g.Stage {
	case STANDBY:
		msg := "Press 'Enter' to start up the workshop!"
		ebitenutil.DebugPrintAt(screen, msg, (ScreenWidth/2)-128, 0)
	case PROCESSING:
		msg := "Processing..."
		ebitenutil.DebugPrintAt(screen, msg, (ScreenWidth/2)-128, 0)
	case COMPLETED:
		msg := "All of the presents have been created!"
		ebitenutil.DebugPrintAt(screen, msg, (ScreenWidth/2)-128, 0)
		msg2 := "Press 'Enter' to start up the workshop!"
		ebitenutil.DebugPrintAt(screen, msg2, (ScreenWidth/2)-128, 14)
	}

	workshopW, _ := WorkshopImg.Size()
	workshopImgOp := &ebiten.DrawImageOptions{}
	workshopImgOp.GeoM.Scale(2, 2)
	workshopImgOp.GeoM.Translate(float64(ScreenWidth-2*workshopW), 0)
	screen.DrawImage(WorkshopImg, workshopImgOp)

	for e := range g.VisElves {
		MVisElves.Lock()
		g.VisElves[e].drawElf(screen)
		MVisElves.Unlock()
	}

	if len(g.VisRoute) > 0 {
		_, h := WorkshopImg.Size()
		x := 0
		y := 2 * h
		msg := " Route:\n "+g.VisRoute
		text.Draw(screen, msg, Font, x, y, color.White)
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
		Size:    12,
		DPI:     dpi,
		Hinting: font.HintingFull,
	})
}