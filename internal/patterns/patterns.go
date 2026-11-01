// Package patterns defines bounded, portable starting grids without graphics or file I/O.
package patterns

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"strings"

	"game-of-life/internal/simulation"
)

const MaxSide = 64
const MaxTextBytes = 8192

//go:embed library.json
var libraryJSON []byte

type Pattern struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Grid        string `json:"grid"`
}

func Library() ([]Pattern, error) {
	var library []Pattern
	if err := json.Unmarshal(libraryJSON, &library); err != nil {
		return nil, err
	}
	return library, nil
}

// Parse accepts rectangular . (dead) / O (alive) rows, with LF or CRLF endings.
// Outer whitespace is ignored; internal spaces, blank rows and other symbols are errors.
func Parse(text string) ([]string, error) {
	if len(text) > MaxTextBytes {
		return nil, fmt.Errorf("starting grid must be at most %d bytes", MaxTextBytes)
	}
	text = strings.TrimSpace(strings.ReplaceAll(text, "\r\n", "\n"))
	rows := strings.Split(text, "\n")
	width := len(rows[0])
	if width < 1 || width > MaxSide || len(rows) > MaxSide {
		return nil, fmt.Errorf("starting grid must be 1..%d columns and rows", MaxSide)
	}
	for y, row := range rows {
		if len(row) != width {
			return nil, fmt.Errorf("starting grid row %d must have %d columns", y+1, width)
		}
		for _, cell := range row {
			if cell != '.' && cell != 'O' {
				return nil, fmt.Errorf("starting grid uses only . (dead) and O (alive), invalid row %d", y+1)
			}
		}
	}
	return rows, nil
}

// Cells centers a selected pattern on an otherwise dead board; it never clips.
func Cells(width, height int, name, custom string) ([]bool, error) {
	if err := (simulation.Config{Width: width, Height: height, CellSize: 1}).Validate(); err != nil {
		return nil, err
	}
	text := custom
	if name != "custom" {
		if custom != "" {
			return nil, fmt.Errorf("custom grid requires pattern custom")
		}
		library, err := Library()
		if err != nil {
			return nil, err
		}
		found := false
		for _, pattern := range library {
			if pattern.ID == name {
				text = pattern.Grid
				found = true
				break
			}
		}
		if !found {
			return nil, fmt.Errorf("unknown starting pattern %q", name)
		}
	}
	rows, err := Parse(text)
	if err != nil {
		return nil, err
	}
	if len(rows) > height || len(rows[0]) > width {
		return nil, fmt.Errorf("starting pattern is %dx%d; increase board dimensions from %dx%d", len(rows[0]), len(rows), width, height)
	}
	cells := make([]bool, width*height)
	x0, y0 := (width-len(rows[0]))/2, (height-len(rows))/2
	for y, row := range rows {
		for x, cell := range row {
			cells[(y0+y)*width+x0+x] = cell == 'O'
		}
	}
	return cells, nil
}
