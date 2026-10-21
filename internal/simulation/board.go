// Package simulation implements Conway's finite, non-wrapping board without graphics.
package simulation

import (
	"fmt"
	"math"
)

// Config describes the initial board and desktop viewport.
type Config struct {
	Width, Height, CellSize, Speed int
	Density                        float64
}

// Validate bounds allocation and viewport dimensions before graphics startup.
func (c Config) Validate() error {
	if c.Width < 1 || c.Height < 1 || c.Width > 4096 || c.Height > 4096 || c.Width*c.Height > 4_000_000 {
		return fmt.Errorf("board dimensions must be 1..4096 with at most 4,000,000 cells")
	}
	if c.CellSize < 1 || c.CellSize > 8192 || c.Width > 8192/c.CellSize || c.Height > 8192/c.CellSize {
		return fmt.Errorf("cell size must be positive and viewport dimensions at most 8192 pixels")
	}
	if c.Speed < 0 || c.Speed > 60 {
		return fmt.Errorf("speed must be 0..60")
	}
	if math.IsNaN(c.Density) || math.IsInf(c.Density, 0) || c.Density < 0 || c.Density > 1 {
		return fmt.Errorf("variance must be a finite probability between 0 and 1")
	}
	return nil
}

// Board owns two simulation buffers and the RGBA image of its current generation.
// Cells and Pixels are read by the renderer; change cells through Reset.
type Board struct {
	Width, Height          int
	Cells                  []bool
	Pixels                 []byte
	Generation, Population int
	buffer                 []bool
	tickCounter            int
}

func New(width, height int, density float64, random func() float64) (*Board, error) {
	if err := (Config{Width: width, Height: height, CellSize: 1, Density: density}).Validate(); err != nil {
		return nil, err
	}
	if random == nil {
		return nil, fmt.Errorf("random source is required")
	}
	b := &Board{Width: width, Height: height, Cells: make([]bool, width*height), buffer: make([]bool, width*height), Pixels: make([]byte, width*height*4)}
	for i := range b.Cells {
		b.Cells[i] = random() < density
	}
	b.refresh()
	return b, nil
}

// Reset copies a complete starting pattern and resets generation and tick timing.
func (b *Board) Reset(cells []bool) error {
	if len(cells) != len(b.Cells) {
		return fmt.Errorf("pattern has %d cells; want %d", len(cells), len(b.Cells))
	}
	copy(b.Cells, cells)
	clear(b.buffer)
	b.Generation = 0
	b.tickCounter = 0
	b.refresh()
	return nil
}

// Advance preserves the original game's strict tick threshold, including pause.
// Speed is a nominal setting, not an exact generation-per-second guarantee.
func (b *Board) Advance(speed, tps int) bool {
	if speed != 0 && float32(b.tickCounter) > float32(tps)/float32(speed) {
		b.Step()
		b.tickCounter = 0
		return true
	}
	b.tickCounter++
	return false
}

// Step applies all births and deaths simultaneously; neighbors outside the board are dead.
func (b *Board) Step() {
	for y := range b.Height {
		for x := range b.Width {
			neighbors := 0
			for dy := -1; dy <= 1; dy++ {
				for dx := -1; dx <= 1; dx++ {
					if dx == 0 && dy == 0 {
						continue
					}
					nx, ny := x+dx, y+dy
					if nx >= 0 && ny >= 0 && nx < b.Width && ny < b.Height && b.Cells[ny*b.Width+nx] {
						neighbors++
					}
				}
			}
			i := y*b.Width + x
			b.buffer[i] = neighbors == 3 || (b.Cells[i] && neighbors == 2)
		}
	}
	b.Cells, b.buffer = b.buffer, b.Cells
	b.Generation++
	b.refresh()
}

func (b *Board) refresh() {
	b.Population = 0
	clear(b.Pixels)
	for i, alive := range b.Cells {
		if alive {
			b.Population++
			b.Pixels[4*i+1] = 255
			b.Pixels[4*i+3] = 255
		}
	}
}
