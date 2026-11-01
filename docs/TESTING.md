# Testing

Run from the repository root with the Go version declared in go.mod:

```sh
go test ./...
go vet ./...
go mod verify
go mod tidy -diff
go build -o build/game .
go test -race -count=1 -coverprofile=build/coverage.out ./internal/simulation ./internal/patterns
go tool cover -func=build/coverage.out
go tool cover -html=build/coverage.out -o build/coverage.html
```

Create build/ first if it does not exist. In PowerShell, quote arguments containing equals signs or dots, for example '-coverprofile=build/coverage.out'. The race detector requires a supported platform and C toolchain; regular go test works without that additional toolchain for the pure simulation package.

The simulation package has no display, graphics, audio or filesystem dependency. CI runs its tests with the race detector on Linux and publishes raw/HTML coverage as an artifact and a summary. A Windows job also runs all package tests, vet, module verification and native/Wasm builds. Simulation coverage does not include the Ebitengine adapter, keyboard events, render output or asset loading.

WebAssembly compilation (POSIX shell):

```sh
GOOS=js GOARCH=wasm go build -o build/game.wasm .
```

PowerShell:

```powershell
$env:GOOS = 'js'
$env:GOARCH = 'wasm'
try { go build -o build/game.wasm . } finally { Remove-Item Env:GOOS; Remove-Item Env:GOARCH }
```

This is the low-level compilation check; the complete browser launcher and browser-specific manual checks are documented in BROWSER.md. Native Linux packages importing Ebitengine also require graphics libraries and a display; use Xvfb when testing those packages on a headless Linux host. See the Ebitengine installation guide linked in README.

## Automated behavior

Tests cover block still life, blinker, glider motion, finite boundaries, all 512 possible 3x3 starting boards for three generations, population/RGBA synchronization, deterministic random initialization, reset/copy isolation, per-board timing, pause/resume and startup limits (including NaN/infinite density).

## Manual checklist

Run go run .; confirm random cells and population are visible immediately; try -cell 1 -speed 0 and verify initial pixels remain visible while paused; restart with R; toggle D and L; adjust speed with equals/minus; toggle fullscreen with F and restore; quit with Q. Check -variance 0 and 1; confirm invalid cell size, dimensions, speed and variance return clear errors. The one-pixel fast path does not draw grid lines. Rendering, fullscreen, keyboard input and visual scaling still need a real display check.


## Browser checks

Run `node --test scripts/*.test.mjs`, `go run ./cmd/web`, and `node scripts/package-pages.mjs <fresh-output-directory>`. CI uploads the complete Pages artifact after these pass. See BROWSER.md for the actual Chromium checks performed and remaining browser/touch coverage limits. Verify settings after restart, initial focus, graceful failed loading, gzip bundle startup, and that the preview URL appears on the PR and in the workflow summary.

Install browser assets with `npm ci --ignore-scripts` before browser tests/builds. Shared loader/status/frame behavior tests live in web-shared; this repo tests its own packaged references and options/storage where applicable.

## Starting grids (issue #7)

Go tests cover CRLF/text limits, invalid symbols/ragged grids, centering/non-square
boards, no clipping, empty grids, invalid selections, all built-in entries and
known evolution (still life/oscillator periods and glider/spaceship displacement).
Startup tests cover initial counters/pixels, repeatable named/custom restart,
atomic failure, unchanged random density, bounded file reads and missing files.
The Linux race/coverage job includes internal/patterns with internal/simulation.

Browser tests use the actual shared catalog, cover query encoding/fit checks,
name replacement/deletion/limits, storage getter/read/write/quota/corruption
failures, literal saved names, UI preset-to-custom/save/load/delete/paused wiring,
and library-fetch failure. Packaging checks include the new modules and catalog
in the immutable asset set. Actual Chromium verification is recorded in BROWSER.md.
