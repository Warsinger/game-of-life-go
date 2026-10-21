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
* R restart with a new random board
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

Run `go run ./cmd/web -serve -port 8080` and open [localhost:8080](http://localhost:8080/). The game starts immediately; every command-line setting is also available on the page. Apply / restart creates a fresh run. The default board starts evolving randomly.

The deployment target is [golgo.pages.dev](https://golgo.pages.dev/). See [browser behavior](docs/BROWSER.md) and [deployment/preview instructions](docs/DEPLOYMENT.md). Application icons are tracked separately in #8.
