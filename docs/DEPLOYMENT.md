# Cloudflare Pages

Project: golgo. Expected hostname: https://golgo.pages.dev/. Production branch: main. Use a free Direct Upload Pages project, not Workers or Git-connected Pages builds.

Repository Actions secrets CLOUDFLARE_API_TOKEN (Pages Edit on the account) and CLOUDFLARE_ACCOUNT_ID are already configured by the owner. They are passed explicitly only to the shared deployment workflow; source code and build artifacts contain no credential values.

## Build and publish

The Build workflow checks native compilation, simulation tests/races and browser tests/bundling. A successful main push publishes production automatically. PR runs build/test without secrets or publishing. For a manual preview, run Build on the open same-repository PR branch with deploy=true and preview_pr=<number>. A manual main run with deploy=true republishes the current main revision. The default manual run only builds.

The reusable Warsinger/cloudflare-pages-deploy workflow validates the exact project, hostname, production branch, current main/PR head and artifact revision before upload. It posts or updates a preview link on the PR and includes the URL as a link in the build summary. Stale, closed or fork PRs are rejected. The shared implementation is referenced by its reviewed immutable commit, matching deployment-ref, as required by its caller validation contract. Standard GitHub Actions use major tags per owner preference.

Build the upload directory locally with:

```sh
go run ./cmd/web
node scripts/package-pages.mjs build/pages
```

Choose a fresh output directory for each packaging run. The packager never deletes caller data. It compresses the Wasm binary, hashes the entire runtime/settings/style asset set into an immutable directory, enforces the 25 MiB per-file limit, writes revision metadata, and makes entry HTML/version metadata revalidate. A custom 404 prevents missing asset requests from silently serving the home page.

## Recovery and follow-up

For a code regression, revert the bad commit on main and let the guarded workflow publish the new main revision. Cloudflare dashboard rollback can restore an earlier upload for an urgent visual check, but keep repository history aligned with the intended release; later main pushes publish the latest main again. After a successful production deployment, verify startup/settings/controls and version.json, then merge the companion landing-page catalog update. Application icon sources and platform support are documented in ICONS.md.

The browser job installs the integrity-locked shared browser package with npm ci --ignore-scripts. Shared modules/styles are included in the app asset hash and Pages upload. No hosting target, credentials or production publication policy changes are required.
