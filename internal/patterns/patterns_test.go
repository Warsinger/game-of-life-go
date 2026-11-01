package patterns

import (
	"reflect"
	"strings"
	"testing"

	"game-of-life/internal/simulation"
)

func TestParseAndLimits(t *testing.T) {
	rows, err := Parse("\r\n.O.\r\n..O\r\nOOO\r\n")
	if err != nil || !reflect.DeepEqual(rows, []string{".O.", "..O", "OOO"}) {
		t.Fatalf("CRLF: %v, %v", rows, err)
	}
	for _, text := range []string{"", "O.\nO", "O \nOO", "O\n\nO", "o", "😀", strings.Repeat("O", 65), strings.Repeat("O\n", 65), strings.Repeat("O", MaxTextBytes+1)} {
		if _, err := Parse(text); err == nil {
			t.Errorf("accepted invalid grid %q", text[:min(len(text), 80)])
		}
	}
}

func TestCenteredSeedAndValidation(t *testing.T) {
	cells, err := Cells(7, 5, "glider", "")
	if err != nil {
		t.Fatal(err)
	}
	want := map[int]bool{1*7 + 3: true, 2*7 + 4: true, 3*7 + 2: true, 3*7 + 3: true, 3*7 + 4: true}
	for i, v := range cells {
		if v != want[i] {
			t.Fatalf("cell %d: %v", i, v)
		}
	}
	for _, tc := range []struct {
		w, h       int
		name, text string
	}{{2, 2, "glider", ""}, {5, 5, "missing", ""}, {0, 5, "block", ""}, {4096, 4096, "block", ""}, {5, 5, "glider", "O"}, {5, 5, "custom", ""}} {
		if _, err := Cells(tc.w, tc.h, tc.name, tc.text); err == nil {
			t.Fatalf("accepted invalid selection %+v", tc)
		}
	}
	empty, err := Cells(2, 2, "custom", "..\n..")
	if err != nil {
		t.Fatal(err)
	}
	for _, v := range empty {
		if v {
			t.Fatal("empty grid became live")
		}
	}
}

func TestLibraryPatternsAndKnownEvolution(t *testing.T) {
	library, err := Library()
	if err != nil {
		t.Fatal(err)
	}
	ids := map[string]bool{}
	for _, p := range library {
		if ids[p.ID] || p.ID == "" || p.Name == "" || p.Description == "" {
			t.Fatalf("bad library record %+v", p)
		}
		ids[p.ID] = true
		cells, err := Cells(20, 20, p.ID, "")
		if err != nil {
			t.Fatal(err)
		}
		board, err := simulation.New(20, 20, 0, func() float64 { return 1 })
		if err != nil {
			t.Fatal(err)
		}
		if err := board.Reset(cells); err != nil {
			t.Fatal(err)
		}
		if board.Generation != 0 || board.Population == 0 {
			t.Fatalf("bad initial counters for %s", p.ID)
		}
		steps := map[string]int{"block": 1, "blinker": 2, "beacon": 2, "glider": 4, "lwss": 4}[p.ID]
		for range steps {
			board.Step()
		}
		if p.ID == "block" || p.ID == "blinker" || p.ID == "beacon" {
			if !reflect.DeepEqual(board.Cells, cells) {
				t.Fatalf("%s did not repeat", p.ID)
			}
		}
		if p.ID == "glider" {
			for y := 1; y < 20; y++ {
				for x := 1; x < 20; x++ {
					if board.Cells[y*20+x] != cells[(y-1)*20+x-1] {
						t.Fatal("glider did not travel diagonally after four generations")
					}
				}
			}
		}
		if p.ID == "lwss" {
			for y := 0; y < 20; y++ {
				for x := 0; x < 18; x++ {
					if board.Cells[y*20+x] != cells[y*20+x+2] {
						t.Fatal("spaceship did not travel left after four generations")
					}
				}
			}
		}
	}
	if len(library) < 6 {
		t.Fatal("pattern library incomplete")
	}
}
