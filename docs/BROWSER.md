# Browser application

The page starts immediately with the same defaults as the desktop app and a random board. There is no start screen. All command-line options have labeled web controls. Apply / restart validates the settings and creates a fresh iframe/game instance; R randomizes in place. Board dimensions are fixed for each run.

The game canvas receives focus after its first actual Go update. Use Focus game after interacting with settings, or click the game. Browser Q is ignored; desktop Q still exits. Fullscreen is available from the page button or the game F key when supported by the browser. Grid lines retain the desktop limitation in the one-pixel fast path.

## Local build

Run from the repository root:

```sh
npm ci --ignore-scripts
go run ./cmd/web -serve -port 8080
```

Open http://localhost:8080/. Without -serve, the command writes build/web and exits. Serve over HTTP rather than opening HTML as a local file. The builder copies wasm_exec.js from the exact Go toolchain that compiled the binary. It replaces inherited GOOS/GOARCH values for the Wasm subprocess without changing the native caller.

Settings are passed as allowlisted query arguments and validated both in JavaScript and by Go before game setup. They are not saved between visits. Browser rules and native simulation use the same Go implementation. No multiplayer, persistence, editing interface or new physics is added.

## Loading and security boundaries

The launcher receives same-origin messages only from its own iframe and current session token, discarding previous-run and foreign messages. It reports ready only after Go's first update; download/runtime failures remain visible with a retry path. Runtime failures do not become successful loading messages. A slow-load notice does not prevent late success. Fullscreen denial is recoverable.

The local loader accepts raw Wasm. The Pages packager generates gzip Wasm and the same loader recognizes gzip magic bytes, avoiding double decompression when a browser/server already decodes the response. A current browser with WebAssembly, WebGL and DecompressionStream is required for the hosted bundle.

## Verification

Run node --test scripts/*.test.mjs plus the Go checks in TESTING.md. JavaScript tests cover option conversion/limits, real HTML/runtime references, reproducible hashed assets, gzip byte integrity, asset limits, immutable caching, loader success/failure and input message trust/release behavior. Coverage of actual rendering remains a browser check.

Manually verified in the in-app Chromium browser: automatic startup, keyboard focus, random evolution, rejected oversized allocation, restart and a paused one-pixel all-live board. Review preview deployment on other browsers and real touch devices before claiming a full support matrix.

## Shared browser package

Node 24+ installs the exact integrity-locked Warsinger/web-shared commit before browser builds: npm ci --ignore-scripts. Native builds remain independent of npm. The builder copies the required shared modules/styles into build/web/shared/ and keeps wasm_exec.js from its own compiling Go toolchain. Pages hashes these files with the app, so every deployment is self-contained and makes no cross-site runtime requests. Update shared code through a reviewed dependency/lockfile change and verify a preview.

Life keeps its random immediate startup, settings conversion/validation and Conway simulation. The shared settings launcher handles restart, canvas focus, fullscreen, trusted current-session messages and status. Footer HTML/links remain local; one shared footer stylesheet supports each layout. Common lifecycle/Go-loader tests are in web-shared. Consumer tests check options where applicable and real packaged references, icons, gzip and shared-style invalidation.
