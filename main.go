package main

import (
	"flag"
	"fmt"
	"game-of-life/internal/platform"
	"game-of-life/internal/simulation"
	"image/color"
	"log"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type GameInfo struct {
	*simulation.Board
	width, height, cellSize, speed int
	variance                       float64
	debug, lines                   bool
}

var cellColor = color.RGBA{0, 255, 0, 255}

const minSpeed = 0
const maxSpeed = 60

func (g *GameInfo) Draw(screen *ebiten.Image) {
	screen.Clear()

	size := screen.Bounds().Size()
	if g.cellSize == 1 && len(g.Pixels) == size.X*size.Y*4 {
		// if size of cells is 1 then use pixel method to make images
		screen.WritePixels(g.Pixels)
	} else {
		for x := range g.width {
			for y := range g.height {
				if g.Cells[y*g.width+x] {
					vector.DrawFilledRect(screen, float32(x*g.cellSize), float32(y*g.cellSize), float32(g.cellSize), float32(g.cellSize), cellColor, true)
				}
			}
		}

		if g.lines {
			for i := 0; i <= size.Y; i += g.cellSize {
				vector.StrokeLine(screen, 0, float32(i), float32(size.X), float32(i), 1, color.White, true)
			}
			for i := 0; i <= size.X; i += g.cellSize {
				vector.StrokeLine(screen, float32(i), 0, float32(i), float32(size.Y), 1, color.White, true)
			}
		}
	}

	if g.debug {
		str := fmt.Sprintf("Generation %v\nPopulation %v\nSpeed %v\nTPS %2.1f", g.Generation, g.Population, g.speed, ebiten.ActualTPS())
		ebitenutil.DebugPrintAt(screen, str, 5, 5)
	}
}

func (g *GameInfo) Layout(outsideWidth int, outsideHeight int) (screenWidth int, screenHeight int) {
	if platform.Browser() {
		return g.width * g.cellSize, g.height * g.cellSize
	}
	return outsideWidth, outsideHeight
}

func (g *GameInfo) Update() error {
	platform.Ready()
	if !platform.Browser() && inpututil.IsKeyJustPressed(ebiten.KeyQ) {
		return ebiten.Termination
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyR) {
		return g.Init()
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyD) {
		g.debug = !g.debug
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyL) {
		g.lines = !g.lines
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyF) {
		ebiten.SetFullscreen(!ebiten.IsFullscreen())
		if !ebiten.IsFullscreen() {
			ebiten.SetWindowSize(g.width*g.cellSize, g.height*g.cellSize)
		}
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyEqual) {
		g.speed = (min(g.speed+5, maxSpeed))
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyMinus) {
		g.speed = (max(g.speed-5, minSpeed))
	}

	g.Board.Advance(g.speed, ebiten.TPS())
	return nil
}

func (g *GameInfo) Init() error {
	config := simulation.Config{Width: g.width, Height: g.height, CellSize: g.cellSize, Speed: g.speed, Density: g.variance}
	if err := config.Validate(); err != nil {
		return err
	}
	board, err := simulation.New(g.width, g.height, g.variance, rand.Float64)
	if err != nil {
		return err
	}
	g.Board = board
	return nil
}

func main() {
	width := flag.Int("width", 80, "Board width in cells")
	height := flag.Int("height", 80, "Board height in cells")
	cellSize := flag.Int("cell", 10, "Size of each cell in pixels")
	speed := flag.Int("speed", 15, "Nominal speed, min 0 (pause) max 60, + or - to adjust in game")
	variance := flag.Float64("variance", 0.5, "Initial live-cell probability, from 0 to 1")
	debug := flag.Bool("debug", false, "Show debug info, D to toggle in game")
	lines := flag.Bool("lines", false, "Draw grid lines, L to toggle in game")

	if err := flag.CommandLine.Parse(platform.Args()); err != nil {
		log.Fatal(err)
	}

	g := &GameInfo{width: *width, height: *height, cellSize: *cellSize, variance: *variance, debug: *debug, lines: *lines, speed: *speed}
	if err := g.Init(); err != nil {
		log.Fatal(err)
	}
	ebiten.SetWindowSize(g.width*g.cellSize, g.height*g.cellSize)
	ebiten.SetWindowTitle("Conway's Game of Life")
	ebiten.SetTPS(60)

	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
