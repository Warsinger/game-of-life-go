package main

import (
	"flag"
	"fmt"
	"image/color"
	"log"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type GameInfo struct {
	width      int
	height     int
	cellSize   int
	variance   float64
	cells      []bool
	buffer     []bool
	pixels     []byte
	generation int
	population int
	speed      int
	debug      bool
	lines      bool
}

var cellColor = color.RGBA{0, 255, 0, 255}

const minSpeed = 0
const maxSpeed = 60

var tickCounter = 0

func (g *GameInfo) Draw(screen *ebiten.Image) {
	screen.Clear()

	size := screen.Bounds().Size()
	if g.cellSize == 1 && len(g.pixels) == size.X*size.Y*4 {
		// if size of cells is 1 then use pixel method to make images
		screen.WritePixels(g.pixels)
		// TODO fix bug with writing pixels in full screen mode
	} else {
		for x := range g.width {
			for y := range g.height {
				if g.cells[y*g.width+x] {
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
		str := fmt.Sprintf("Generation %v\nPopulation %v\nSpeed %v\nTPS %2.1f", g.generation, g.population, g.speed, ebiten.ActualTPS())
		ebitenutil.DebugPrintAt(screen, str, 5, 5)
	}
}

func (g *GameInfo) Layout(outsideWidth int, outsideHeight int) (screenWidth int, screenHeight int) {
	return outsideWidth, outsideHeight
}

func (g *GameInfo) Update() error {
	if inpututil.IsKeyJustPressed(ebiten.KeyQ) {
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

	if g.speed != 0 && float32(tickCounter) > float32(ebiten.TPS())/float32(g.speed) {
		g.generation++
		err := g.UpdatePopulation()
		if err != nil {
			return err
		}
		tickCounter = 0
	} else {
		tickCounter++
	}
	return nil
}

func (g *GameInfo) Init() error {
	g.cells = make([]bool, g.width*g.height)
	g.randomData()
	g.buffer = make([]bool, g.width*g.height)
	g.pixels = make([]byte, g.width*g.height*4)
	g.generation = 0

	return nil
}

func (g *GameInfo) randomData() {
	for x := range g.cells {
		if rand.Float64() < g.variance {
			g.cells[x] = true
		}
	}
}

// Any live cell with fewer than two live neighbours dies, as if by underpopulation.
// Any live cell with two or three live neighbours lives on to the next generation.
// Any live cell with more than three live neighbours dies, as if by overpopulation.
// Any dead cell with exactly three live neighbours becomes a live cell, as if by reproduction.

func (g *GameInfo) UpdatePopulation() error {
	g.population = 0
	for y := range g.height {
		for x := range g.width {

			neighbors := g.countNeighbors(x, y)

			idx := y*g.width + x
			if !g.cells[idx] && neighbors == 3 {
				g.buffer[idx] = true
			} else if neighbors < 2 || neighbors > 3 {
				g.buffer[idx] = false
			} else {
				g.buffer[idx] = g.cells[idx]
			}
			if g.buffer[idx] {
				g.population++
				g.pixels[idx*4] = cellColor.R
				g.pixels[idx*4+1] = cellColor.G
				g.pixels[idx*4+2] = cellColor.B
				g.pixels[idx*4+3] = cellColor.A
			} else {
				g.pixels[idx*4] = 0
				g.pixels[idx*4+1] = 0
				g.pixels[idx*4+2] = 0
				g.pixels[idx*4+3] = 0
			}
		}
	}
	temp := g.buffer
	g.buffer = g.cells
	g.cells = temp
	return nil
}

func (g *GameInfo) countNeighbors(x, y int) uint8 {
	var neighbors uint8 = 0
	for k := -1; k <= 1; k++ {
		for m := -1; m <= 1; m++ {
			if k == 0 && m == 0 {
				continue
			}
			x2 := x + m
			y2 := y + k
			if x2 < 0 || y2 < 0 || x2 >= g.width || y2 >= g.height {
				continue
			}
			if g.cells[y2*g.width+x2] {
				neighbors++
			}
		}
	}
	return neighbors
}

func main() {
	width := flag.Int("width", 80, "Board width in cells")
	height := flag.Int("height", 80, "Board height in cells")
	cellSize := flag.Int("cell", 10, "Size of each cell in pixels")
	speed := flag.Int("speed", 15, "Generations per second, min 0 max 60, + or - to adjust in game")
	variance := flag.Float64("variance", 0.5, "Grid Size")
	debug := flag.Bool("debug", false, "Show debug info, D to toggle in game")
	lines := flag.Bool("lines", false, "Draw grid lines, L to toggle in game")

	flag.Parse()

	g := &GameInfo{width: *width, height: *height, cellSize: *cellSize, variance: *variance, debug: *debug, lines: *lines, speed: *speed}
	g.Init()
	ebiten.SetWindowSize(g.width*g.cellSize, g.height*g.cellSize)
	ebiten.SetWindowTitle("Conway's Game of Life")
	ebiten.SetTPS(60)

	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
