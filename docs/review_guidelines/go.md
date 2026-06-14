# Go Review Guidelines

Applies to: `**/*.go`

- Keep functions short and focused. ~30 lines is a directional
  guideline, not a hard limit — test helpers and generated code
  are exempt.
- Each function should do one thing. Prefer named helper functions
  over inline anonymous functions, but short idiomatic closures
  (e.g. `go func` callbacks under ~10 lines) are fine.
- Flag functions that are long or complex enough to harm
  readability for someone unfamiliar with the codebase.
