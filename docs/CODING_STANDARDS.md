# Coding standards

Use idiomatic Go, gofmt and small functions. Prefer the standard library; explain new dependencies in the PR. Keep the module path stable unless a migration requires changing it.

Simulation owns mutable state for one instance. The Ebitengine adapter reads state for drawing and forwards input and tick updates. Do not import graphics or audio packages into internal/simulation. Tests must be deterministic; inject a random function rather than reseeding global randomness.

Validate user configuration before allocating large buffers or loading assets. Return useful errors to callers. Avoid silent coercion of invalid startup options. Document numeric bounds and behavior changes. Preserve existing movement/rules during structural refactors.

Tests should assert observable rules and known trajectories, including boundaries and error cases, rather than copying the production implementation. Add regression tests when fixing a bug. Coverage indicates exercised statements, not graphical correctness.

Keep assets embedded and retain their attribution. Never commit generated binaries, local helper copies, secrets or temporary output. Pin GitHub Actions to full commit SHAs with version comments and least-privilege permissions.
