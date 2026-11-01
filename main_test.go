package main

import (
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"game-of-life/internal/patterns"
)

func TestStartingPatternInitializationAndRestart(t *testing.T) {
	for _, name := range []string{"block", "glider", "custom"} {
		t.Run(name, func(t *testing.T) {
			g := &GameInfo{width: 9, height: 7, cellSize: 10, speed: 15, variance: .5, pattern: name}
			if name == "custom" {
				g.seed = ".O.\n..O\nOOO"
			}
			if err := g.Init(); err != nil {
				t.Fatal(err)
			}
			original := append([]bool(nil), g.Cells...)
			if g.Generation != 0 || g.Population == 0 {
				t.Fatal("starting counters missing")
			}
			for i, alive := range g.Cells {
				if (g.Pixels[4*i+3] == 255) != alive {
					t.Fatal("starting pixels out of sync")
				}
			}
			g.Step()
			g.Step()
			if err := g.Init(); err != nil {
				t.Fatal(err)
			}
			if g.Generation != 0 || !reflect.DeepEqual(g.Cells, original) {
				t.Fatal("restart did not restore starting grid")
			}
			previous := g.Board
			g.pattern = "custom"
			g.seed = "bad"
			if err := g.Init(); err == nil || g.Board != previous {
				t.Fatal("invalid seed replaced the active board")
			}
		})
	}
}
func TestDefaultRandomModeStillUsesDensity(t *testing.T) {
	for _, density := range []float64{0, 1} {
		g := &GameInfo{width: 3, height: 4, cellSize: 10, variance: density}
		if err := g.Init(); err != nil {
			t.Fatal(err)
		}
		if g.Population != int(density*12) {
			t.Fatal("random mode density changed")
		}
	}
	g := &GameInfo{width: 3, height: 3, cellSize: 10, pattern: "random", seed: "O"}
	if err := g.Init(); err == nil {
		t.Fatal("random unexpectedly accepted custom text")
	}
}
func TestPatternFileLimitsAndErrors(t *testing.T) {
	name := filepath.Join(t.TempDir(), "seed.cells")
	if _, err := readPatternFile(name); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("missing file: %v", err)
	}
	if err := os.WriteFile(name, []byte("OO\nOO\n"), 0600); err != nil {
		t.Fatal(err)
	}
	text, err := readPatternFile(name)
	if err != nil || text != "OO\nOO\n" {
		t.Fatalf("pattern file %q: %v", text, err)
	}
	if err := os.WriteFile(name, []byte(strings.Repeat("O", patterns.MaxTextBytes+1)), 0600); err != nil {
		t.Fatal(err)
	}
	if _, err := readPatternFile(name); err == nil {
		t.Fatal("unbounded pattern file accepted")
	}
}
