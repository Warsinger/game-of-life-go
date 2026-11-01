# Project specification

## Purpose and current scope

A desktop and browser Conway's Game of Life sandbox written in Go/Ebitengine. A finite rectangular board has no edge wrapping: neighbors outside the board are dead. Live cells survive with two or three neighbors; dead cells are born with exactly three. All changes occur simultaneously.

Random mode is the immediate-start default; initialization and R restart generate an independent board using the configured live-cell probability. Named/custom modes instead center a deterministic starting grid on dead cells, and R restores that selection. Generation starts at zero; population and RGBA pixels are coherent immediately, and restart resets tick timing. The browser offers a plain-text grid editor and an explicitly saved local library of starting patterns, not running-state persistence. There is no multiplayer.

## Architecture

- main.go: command-line options, Ebitengine lifecycle, keyboard input and drawing.
- internal/simulation: validated configuration, random initialization, double-buffered rules, generation/population statistics, RGBA pixels, reset and tick scheduling.
- internal/patterns: shared embedded pattern catalog, bounded rectangular text parsing and centered placement; no graphics/file I/O. Native file reading stays in main.
- web/patterns.mjs and pattern-ui.mjs: app-specific validation/editor/local saved library; shared shell/lifecycle remain unchanged.
- docs/: specification, coding and testing guidance.
- .github/workflows/build.yml: native/Wasm verification and headless simulation coverage.

The simulation accepts an injected random function for deterministic tests. Its board owns all buffers and tick state; no state is shared between instances. The renderer reads its current cells and pixels.

## Configuration and timing

Defaults: width 80 cells, height 80 cells, cell size 10 pixels, variance 0.5, speed 15, debug/grid off. Dimensions must be 1..4096 with at most 4,000,000 total cells; cell size must be positive and the initial viewport at most 8192 pixels in either direction. Density must be finite and between 0 and 1; speed must be 0..60. Invalid input fails before board allocation/window setup.

Speed is a nominal setting: the original scheduler steps when its counter is strictly greater than TPS/speed and resets the counter afterward. This is not an exact generations-per-second rate. Zero pauses generation updates, but the counter continues accumulating; resuming may advance immediately. This PR intentionally preserves that behavior. R restart resets the counter for the new board.

## Rendering and controls

Cells are green on black. One-pixel cells use a packed RGBA fast path when viewport and board dimensions match; other cases draw rectangles. The fast path does not render grid lines. Fullscreen changes the viewport, not simulation dimensions. Controls are documented in README.

## Follow-up

Issue #1 tracks the browser launcher, web controls for all command-line options, Cloudflare deployment and landing-page integration. A Wasm build alone does not complete that issue. Timing accuracy and one-pixel grid rendering can be evaluated separately without changing Conway's rules.

## Web implementation

The browser version uses the same simulation and drawing code with platform-specific query arguments, ready reporting and fixed logical board dimensions. The page launches immediately. All desktop flag values are exposed in validated settings; restart creates a fresh isolated iframe runtime. Rendering scales to the page while preserving the logical configuration. Browser Q is ignored. See BROWSER.md and DEPLOYMENT.md for implemented hosting and loader details; production verification and landing-page publication complete issue #1.

The browser shell uses a pinned, locally bundled Warsinger/web-shared revision. App-specific content, settings/controls and simulation/storage stay in this repository; native behavior is unchanged. See BROWSER.md for installation and update steps.
