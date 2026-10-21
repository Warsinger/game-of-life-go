# Project workflow

- Read docs/PROJECT_SPEC.md, docs/CODING_STANDARDS.md and docs/TESTING.md before relevant changes.
- Use dedicated codex/ branches and PRs. For stacked work, branch from the prerequisite update and identify the dependency in the PR. Target main for automatic issue linkage; dependency changes disappear from the diff after their prerequisite merges.
- Assign every new issue and PR to Warsinger; request Warsinger review. Apply appropriate labels. Apply priority only when the owner sets it.
- Use Closes #N for completed issues and verify closingIssuesReferences. Keep unfinished work open. Never merge without user authorization.
- Preserve controls, simulation rules, jump physics and timing unless a behavior change is agreed or a documented bug is fixed. Add meaningful regression tests for fixes.
- Inspect CI results before handoff; report graphics/audio checks that were not performed.

# GitHub identity

Use Warsinger Code Builder for all GitHub interactions, not connected GitHub tools or default gh authentication. Read C:/Users/eelwe/.codex/github-identities/README.md first. Select Warsinger/game-of-life-go explicitly both in the helper and gh arguments. If the shared helper does not support this repository, use a temporary repository-scoped compatibility copy under ignored build/agent-work, with identity home C:/Users/eelwe/.codex/github-identities; do not modify the shared access allowlist without authorization. Never print credentials, commit helper credentials, persist tokens or put tokens in remote URLs. Publish commits with builder attribution.

# Development

- Use major-version tags for GitHub Actions (for example actions/checkout@v7 and actions/setup-go@v7), rather than commit SHAs, per the owner preference.
- Target the Go version in go.mod; run gofmt, relevant tests, vet and native/Wasm builds.
- Keep internal/simulation independent of Ebitengine, input, file I/O and graphics/audio initialization. Inject randomness where needed.
- Keep scratch files, compatibility helpers, coverage and binaries under ignored build/.
- Web controls, hosting and landing-page integration belong to issue #1. Use the shared Warsinger/cloudflare-pages-deploy workflow when implementing deployment; keep secrets in GitHub Actions.
- Update README and docs when behavior, configuration, architecture or verification changes.
