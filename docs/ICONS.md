# Application icons

The icon is a deterministic green (`#00ff00`) glider on black, matching the simulation pixels. The editable five-by-five live-cell pattern is `assets/icons/source.json`; this is a game-grid illustration, not a screenshot of a random board. No external artwork or font is used. Cells are rendered directly at each output size to preserve crisp edges.

## Regenerate

Run from the repository root:

```sh
go run ./cmd/icons
go run github.com/akavel/rsrc@v0.10.2 -arch amd64 -ico assets/icons/favicon.ico -o resource_windows_amd64.syso
```

Generated PNG sizes are 16, 32, 48, 64, 128, 180, 192, 256, 512 and 1024 pixels. The ICO contains six sizes through 256; app.icns contains sizes through 1024. Generated icons and the Windows resource are deliberately checked in so normal builds do not require generation tools. This is an exception to the rule against committing local build outputs; executables remain ignored.

## Web and desktop

Both the launcher and standalone game declare the favicon, 192/512 PNGs and 180-pixel Apple touch icon. The web builder copies them into its output. Pages packaging puts declared icons in the content-hashed asset directory, so changes get new URLs; /favicon.ico is also available as the conventional fallback. These icons do not add installable-PWA behavior.

Native windows use embedded 32, 48 and 256 pixel icons through ebiten.SetWindowIcon. The root resource_windows_amd64.syso embeds the ICO into Windows AMD64 executables automatically. Its name excludes other operating systems/architectures; it belongs beside the root main package. Other Windows architectures would need their own generated resource.

The current workflow builds a Windows executable and a browser artifact, but has no macOS app-bundle or Linux desktop-entry packaging. app.icns and icon-256.png are supplied for that future work; a bare macOS executable does not acquire a Finder/Dock icon from ICNS alone, and Linux launcher integration depends on a desktop entry.

## Verification

Run go test ./..., go vet ./..., go run ./cmd/web, and node --test scripts/*.test.mjs. Icon tests decode the native images and every ICO/ICNS representation. Packaging tests verify icon references from both entry pages, the conventional favicon, and cache invalidation when an icon changes. CI regenerates the checked-in resources to catch stale derived files.

Artwork was inspected at 16–256 pixels on light and dark surrounds. Native/Windows-resource and Wasm builds pass; actual Explorer, taskbar, macOS Dock and real-device Apple touch-icon display remain manual checks. There are no gameplay or control changes.
