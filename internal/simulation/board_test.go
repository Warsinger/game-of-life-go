package simulation

import (
	"math"
	"math/rand"
	"slices"
	"testing"
)

func board(t *testing.T, w, h int, cells []bool) *Board {
	t.Helper()
	b, err := New(w, h, 0, func() float64 { return 0.5 })
	if err != nil {
		t.Fatal(err)
	}
	if err := b.Reset(cells); err != nil {
		t.Fatal(err)
	}
	return b
}
func pattern(w, h int, coords ...[2]int) []bool {
	cells := make([]bool, w*h)
	for _, p := range coords {
		cells[p[1]*w+p[0]] = true
	}
	return cells
}
func checkPixels(t *testing.T, b *Board) {
	t.Helper()
	population := 0
	for i, alive := range b.Cells {
		want := []byte{0, 0, 0, 0}
		if alive {
			population++
			want = []byte{0, 255, 0, 255}
		}
		if !slices.Equal(b.Pixels[4*i:4*i+4], want) {
			t.Fatalf("pixel %d = %v, want %v", i, b.Pixels[4*i:4*i+4], want)
		}
	}
	if b.Population != population {
		t.Fatalf("population = %d, want %d", b.Population, population)
	}
}
func TestKnownPatterns(t *testing.T) {
	tests := []struct {
		name          string
		initial, want []bool
		steps         int
	}{
		{"block", pattern(6, 6, [2]int{2, 2}, [2]int{3, 2}, [2]int{2, 3}, [2]int{3, 3}), pattern(6, 6, [2]int{2, 2}, [2]int{3, 2}, [2]int{2, 3}, [2]int{3, 3}), 5},
		{"blinker", pattern(6, 6, [2]int{1, 2}, [2]int{2, 2}, [2]int{3, 2}), pattern(6, 6, [2]int{2, 1}, [2]int{2, 2}, [2]int{2, 3}), 1},
		{"blinker returns", pattern(6, 6, [2]int{1, 2}, [2]int{2, 2}, [2]int{3, 2}), pattern(6, 6, [2]int{1, 2}, [2]int{2, 2}, [2]int{3, 2}), 2},
		{"glider", pattern(6, 6, [2]int{1, 0}, [2]int{2, 1}, [2]int{0, 2}, [2]int{1, 2}, [2]int{2, 2}), pattern(6, 6, [2]int{2, 1}, [2]int{3, 2}, [2]int{1, 3}, [2]int{2, 3}, [2]int{3, 3}), 4},
		{"no wrapping", pattern(6, 6, [2]int{0, 0}, [2]int{5, 0}, [2]int{0, 5}), make([]bool, 36), 1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := board(t, 6, 6, tt.initial)
			for range tt.steps {
				b.Step()
				checkPixels(t, b)
			}
			if !slices.Equal(b.Cells, tt.want) {
				t.Fatalf("cells = %v, want %v", b.Cells, tt.want)
			}
			if b.Generation != tt.steps {
				t.Fatalf("generation = %d", b.Generation)
			}
		})
	}
}

// Enumerate every 3x3 starting board against a simple live-coordinate reference.
// This covers every neighbor count, edge and corner, and simultaneous updates.
func TestAllSmallBoards(t *testing.T) {
	for mask := 0; mask < 512; mask++ {
		cells := make([]bool, 9)
		for i := range cells {
			cells[i] = mask&(1<<i) != 0
		}
		b := board(t, 3, 3, cells)
		for generation := 0; generation < 3; generation++ {
			want := make([]bool, 9)
			for i, alive := range cells {
				neighbors := 0
				for j, other := range cells {
					if other && i != j && abs(i%3-j%3) <= 1 && abs(i/3-j/3) <= 1 {
						neighbors++
					}
				}
				want[i] = neighbors == 3 || (alive && neighbors == 2)
			}
			b.Step()
			checkPixels(t, b)
			if !slices.Equal(b.Cells, want) {
				t.Fatalf("mask %d generation %d: got %v want %v", mask, generation+1, b.Cells, want)
			}
			cells = want
		}
	}
}
func abs(n int) int {
	if n < 0 {
		return -n
	}
	return n
}

func TestInitializationAndReset(t *testing.T) {
	for _, density := range []float64{0, 1} {
		b, err := New(3, 2, density, func() float64 { return 0.5 })
		if err != nil {
			t.Fatal(err)
		}
		if b.Generation != 0 || b.Population != int(density*6) {
			t.Fatalf("initial stats: %+v", b)
		}
		checkPixels(t, b)
	}
	rng1, rng2 := rand.New(rand.NewSource(12)), rand.New(rand.NewSource(12))
	b, err := New(8, 8, 0.5, rng1.Float64)
	if err != nil {
		t.Fatal(err)
	}
	other, err := New(8, 8, 0.5, rng2.Float64)
	if err != nil {
		t.Fatal(err)
	}
	if !slices.Equal(b.Cells, other.Cells) {
		t.Fatal("seeded initialization differs")
	}
	checkPixels(t, b)
	b.Step()
	b.Advance(15, 60)
	original := pattern(8, 8, [2]int{0, 0})
	if err := b.Reset(original); err != nil {
		t.Fatal(err)
	}
	original[0] = false
	if !b.Cells[0] || b.Generation != 0 || b.tickCounter != 0 {
		t.Fatal("reset did not copy and reset state")
	}
	checkPixels(t, b)
	before := slices.Clone(b.Cells)
	if err := b.Reset([]bool{true}); err == nil {
		t.Fatal("accepted incomplete pattern")
	}
	if !slices.Equal(before, b.Cells) {
		t.Fatal("invalid reset changed cells")
	}
}
func TestTickTimingPauseAndInstanceIsolation(t *testing.T) {
	a := board(t, 1, 1, []bool{true})
	b := board(t, 1, 1, []bool{true})
	for range 5 {
		if a.Advance(15, 60) {
			t.Fatal("advanced before strict threshold")
		}
	}
	if !a.Advance(15, 60) || a.Generation != 1 {
		t.Fatal("expected sixth update to step")
	}
	if b.Generation != 0 || b.tickCounter != 0 {
		t.Fatal("boards share timing")
	}
	for range 100 {
		if b.Advance(0, 60) {
			t.Fatal("paused board advanced")
		}
	}
	if b.Generation != 0 {
		t.Fatal("pause changed generation")
	}
	if !b.Advance(15, 60) {
		t.Fatal("legacy accumulated pause timing changed")
	}
}
func TestConfigValidation(t *testing.T) {
	valid := Config{Width: 80, Height: 80, CellSize: 10, Speed: 15, Density: 0.5}
	if err := valid.Validate(); err != nil {
		t.Fatal(err)
	}
	for _, change := range []func(*Config){
		func(c *Config) { c.Width = 0 }, func(c *Config) { c.Height = -1 }, func(c *Config) { c.Width = 4097 }, func(c *Config) { c.Width = 3000; c.Height = 3000; c.CellSize = 1 },
		func(c *Config) { c.CellSize = 0 }, func(c *Config) { c.CellSize = 8193 }, func(c *Config) { c.CellSize = 103 }, func(c *Config) { c.Speed = -1 }, func(c *Config) { c.Speed = 61 },
		func(c *Config) { c.Density = -0.1 }, func(c *Config) { c.Density = 1.1 }, func(c *Config) { c.Density = math.NaN() }, func(c *Config) { c.Density = math.Inf(1) },
	} {
		c := valid
		change(&c)
		if err := c.Validate(); err == nil {
			t.Errorf("accepted %+v", c)
		}
	}
	for _, c := range []Config{{1, 1, 1, 0, 0}, {4096, 1, 2, 60, 1}, {2000, 2000, 1, 1, 0.5}} {
		if err := c.Validate(); err != nil {
			t.Errorf("rejected %+v: %v", c, err)
		}
	}
	if _, err := New(0, 2, 0.5, func() float64 { return 0.5 }); err == nil {
		t.Fatal("accepted invalid board")
	}
	if _, err := New(1, 1, 0.5, nil); err == nil {
		t.Fatal("accepted nil random source")
	}
}
