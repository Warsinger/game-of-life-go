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
	cells      [][]uint8
	buffer     [][]uint8
	generation int
	speed      int
	debug      bool
	lines      bool
}

var cellColor = color.RGBA{0, 255, 0, 255}

// var cellColor = color.RGBA{255, 230, 120, 255}

const minSpeed = 0
const maxSpeed = 60

var tickCounter = 0

func (g *GameInfo) Draw(screen *ebiten.Image) {
	for x := range g.width {
		for y := range g.height {
			if g.cells[x][y] == 1 {
				vector.DrawFilledRect(screen, float32(x*g.cellSize), float32(y*g.cellSize), float32(g.cellSize), float32(g.cellSize), cellColor, true)
			}
		}
	}

	if g.lines {
		size := screen.Bounds().Size()

		for i := 0; i <= size.Y; i += g.cellSize {
			vector.StrokeLine(screen, 0, float32(i), float32(size.X), float32(i), 1, color.White, true)
		}
		for i := 0; i <= size.X; i += g.cellSize {
			vector.StrokeLine(screen, float32(i), 0, float32(i), float32(size.Y), 1, color.White, true)
		}
	}

	if g.debug {
		str := fmt.Sprintf("Generation %v\nPopulation %v\nSpeed %v", g.generation, g.CountPopulation(), g.speed)
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
	g.cells = g.initArray(true)
	g.buffer = g.initArray(false)
	g.generation = 0

	return nil
}

func (g *GameInfo) initArray(randomize bool) [][]uint8 {
	grid := make([][]uint8, g.width)
	for x := range grid {
		grid[x] = make([]uint8, g.height)
		if randomize {
			for y := range g.height {
				if rand.Float64() < g.variance {
					grid[x][y] = 1
				}
			}
		}
	}
	return grid
}

func (g *GameInfo) CountPopulation() uint {
	var count uint = 0
	for x := range g.width {
		for y := range g.height {
			count += uint(g.cells[x][y])
		}
	}
	return count
}

// Any live cell with fewer than two live neighbours dies, as if by underpopulation.
// Any live cell with two or three live neighbours lives on to the next generation.
// Any live cell with more than three live neighbours dies, as if by overpopulation.
// Any dead cell with exactly three live neighbours becomes a live cell, as if by reproduction.

func (g *GameInfo) UpdatePopulation() error {
	for x := range g.width {

		for y := range g.height {
			g.buffer[x][y] = 0

			var neighbors uint8 = 0
			for k := -1; k <= 1; k++ {
				for m := -1; m <= 1; m++ {
					if x+k >= 0 && x+k < g.width && y+m >= 0 && y+m < g.height {
						neighbors += g.cells[x+k][y+m]
					}
				}
			}
			neighbors -= g.cells[x][y]

			if g.cells[x][y] == 0 && neighbors == 3 {
				g.buffer[x][y] = 1
			} else if neighbors < 2 || neighbors > 3 {
				g.buffer[x][y] = 0
			} else {
				g.buffer[x][y] = g.cells[x][y]
			}
		}

	}
	temp := g.buffer
	g.buffer = g.cells
	g.cells = temp
	return nil
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
	fmt.Println(g)
	g.Init()
	ebiten.SetWindowSize(g.width*g.cellSize, g.height*g.cellSize)
	ebiten.SetTPS(60)
	ebiten.SetFullscreen(true)

	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
