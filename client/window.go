package client

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
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
)

const (
	ScreenWidth  = 640
	ScreenHeight = 480
)

// VisElf - stores the id and the position of an elf sprite
type VisElf struct {
	Id     int
	X      float64
	Y      float64
	StartX float64
	StartY float64
	EndX float64
	EndY float64
	Img *ebiten.Image
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
	VisElves []VisElf
	VisQueue chan Task
	WorkshopSpace []int
}

// Update proceeds the game state.
// Update is called every tick (1/60 [s] by default).
func (g *Game) Update() error {
	// Write your game's logical update.
	select {
	case task := <-g.VisQueue: // upon receiving a task, update visualisation
		switch task.Action {
		case util.ELF_ENTER:
			workshopW, workshopH := WorkshopImg.Size()
			mWorkshop.Lock()
			for s := range g.WorkshopSpace {
				if g.WorkshopSpace[s] == -1 {
					g.WorkshopSpace[s] = task.Id
					MVisElves.Lock()
					if (s+1) % 2 == 0 { // if even, go to right
						g.VisElves[task.Id].X = float64(ScreenWidth-workshopW) // half of workshopW size (don't forget it is 2x scaled)
					} else {
						g.VisElves[task.Id].X = float64(ScreenWidth-2*workshopW)
					}
					if s+1 > 2 { // if > 2, go to bottom
						g.VisElves[task.Id].Y = float64(workshopH) // half of workshopH size (don't forget it is 2x scaled)
					} else {
						g.VisElves[task.Id].Y = float64(0)
					}
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
				}
			}
			mWorkshop.Unlock()
			MVisElves.Lock()
			g.VisElves[task.Id].X = g.VisElves[task.Id].StartX
			g.VisElves[task.Id].Y = g.VisElves[task.Id].StartY
			MVisElves.Unlock()
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

	workshopW, _ := WorkshopImg.Size()
	workshopImgOp := &ebiten.DrawImageOptions{}
	workshopImgOp.GeoM.Scale(2, 2)
	workshopImgOp.GeoM.Translate(float64(ScreenWidth-2*workshopW), 0)
	screen.DrawImage(WorkshopImg, workshopImgOp)

	for i := 0; i < 8; i++ {
		MVisElves.Lock()
		elf := g.VisElves[i]
		elf.drawElf(screen)
		MVisElves.Unlock()
	}
}

// Layout takes the outside size (e.g., the window size) and returns the (logical) screen size.
// If you don't have to adjust the screen size with the outside size, just return a fixed size.
func (g *Game) Layout(outsideWidth, outsideHeight int) (sw, sh int) {
	return ScreenWidth, ScreenHeight
}

// DecodeImages - decodes the santa and workshop images into ebiten.image types and stores them in the respective variables
func DecodeImages() {
	var err error
	SantaImg, _, err = ebitenutil.NewImageFromFile("sprites/santa.png")
	util.Check(err)

	ElfImg, _, err = ebitenutil.NewImageFromFile("sprites/elf_0.png") // used only for the size
	util.Check(err)

	WorkshopImg, _, err = ebitenutil.NewImageFromFile("sprites/workshop.png")
	util.Check(err)
}
