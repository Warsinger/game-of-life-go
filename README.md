# Conway's Game of Life
Conway's Game of Life written in Golang for fun and learning

## Rules
* Any live cell with fewer than two live neighbours dies, as if by underpopulation.
* Any live cell with two or three live neighbours lives on to the next generation.
* Any live cell with more than three live neighbours dies, as if by overpopulation.
* Any dead cell with exactly three live neighbours becomes a live cell, as if by reproduction.

## Controls
* Equals/plus or minus to adjust nominal game speed (0 pauses)
* D toggle debug information
* R restore the selected starting grid, or create a new board in Random mode
* L toggle grid lines (except in the one-pixel fast path)
* F toggle fullscreen
* Q quit

## Command line args
* -cell int
  * Size of each cell in pixels (default 10)
* -debug
  * Show debug info, D to toggle in game
* -height int
  * Board height in cells (default 80)
* -lines
  * Draw grid lines, L to toggle in game
* -speed int
  * Nominal speed, min 0 (pause) max 60, + or - to adjust in game (default 15)
* -variance float
  * Initial live-cell probability from 0 to 1 (default 0.5)
* -width int
  * Board width in cells (default 80)
* -pattern string
  * random (default), block, blinker, beacon, glider, r-pentomino, lwss, or custom
* -seed string
  * Newline-separated . / O rows for -pattern custom
* -pattern-file string
  * Native only: read a custom starting grid from a local text file

## Starting patterns

Choose a built-in pattern or edit a rectangular grid using `O` for live cells
and `.` for dead cells. Patterns are centered on an otherwise empty board;
they must fit without clipping. Custom grids allow 1–64 rows/columns (8192-byte
input limit). Outer whitespace and CRLF line endings are accepted; internal
spaces, ragged/blank rows and other symbols are rejected. An all-dead grid is valid.

```sh
go run . -pattern glider -width 40 -height 30 -cell 16 -speed 0
go run . -pattern-file starting-grid.txt -speed 0
```

For example, starting-grid.txt can contain:

```text
.O.
..O
OOO
```

Use `-pattern-file` by itself, without another pattern/seed selection. Files are
read once on launch; R restores that loaded grid, while relaunching rereads the
file. Random mode remains the default, and variance applies only in Random mode.

In the browser, select a preset and edit its grid to make a custom one. **Load
paused** sets speed to zero so you can inspect generation 0; use + in the game
or set a nonzero speed and Apply to start evolution. **My saved patterns** stores
up to 20 named starting grids in this browser, with a session-only fallback when
storage is blocked/full. Reusing a name replaces that saved pattern. Load brings
it into the editor; Apply / restart or Load paused starts it. Saved grids are
portable shapes, not a history of the running simulation, and are not synced
between devices. Page reload still starts a fresh random board; your named library
remains available. No account, backend or cloud storage is used.

See [starting-pattern screenshot](docs/starting-patterns.jpg). The built-ins are
standard Conway patterns: block, blinker, beacon, glider, R-pentomino and lightweight
spaceship; see [LifeWiki](https://conwaylife.com/wiki/Category:Patterns) for background.

## Build and dependency verification

Requires Go 1.27.3 or newer. Go's automatic toolchain selection can download the version declared in `go.mod`.

```sh
go mod download
go mod verify
go vet ./...
go test ./...
go build -o build/game .
```

Ebitengine's native Linux build requires its graphics development libraries and a display to run. See the [Ebitengine installation guide](https://ebitengine.org/en/documents/install.html).

A WebAssembly compilation check is also available:

```sh
GOOS=js GOARCH=wasm go build -o build/game.wasm .
```

In PowerShell, set `$env:GOOS = "js"` and `$env:GOARCH = "wasm"` before the build, then remove those environment variables afterward. This is the lower-level compilation check; the complete browser launcher is built with the command below. Deterministic simulation tests run without a display. See [testing and manual checks](docs/TESTING.md), [project specification](docs/PROJECT_SPEC.md), [coding standards](docs/CODING_STANDARDS.md), and [agent workflow](AGENTS.md).

Board dimensions must be 1..4096 with at most 4,000,000 cells. Cell size must be positive and the initial viewport must fit within 8192 pixels in each dimension. Variance must be finite and in [0, 1]; speed must be 0..60. The historical tick scheduler makes speed approximate.

Run the desktop app with `go run .`, or for example `go run . -width 100 -height 60 -cell 8 -variance 0.3`. Population and cells are initialized before the first step, including when paused.

## Browser version

Run `go run ./cmd/web -serve -port 8080` and open [localhost:8080](http://localhost:8080/). The game starts immediately; simulation settings and the starting-grid editor are available on the page (local pattern-file input is native only). Apply / restart creates a fresh run. The default board starts evolving randomly.

The deployment target is [golgo.pages.dev](https://golgo.pages.dev/). See [browser behavior](docs/BROWSER.md) and [deployment/preview instructions](docs/DEPLOYMENT.md). See [application icons](docs/ICONS.md) for artwork, regeneration and desktop support.

Browser builds now require Node 24+ and `npm ci --ignore-scripts` before the existing Go web command. Shared styling and browser lifecycle come from the exact locked [web-shared](https://github.com/Warsinger/web-shared) revision; native builds and game rules are unchanged.
