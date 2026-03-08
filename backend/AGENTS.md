# Backend Conventions

## Scope

This backend is intentionally small. Optimize for clear package boundaries, explicit dependencies, and readability over framework-heavy abstractions.

## Architecture

- Keep binaries under `cmd/` thin. They should wire config, logging, and app startup only.
- Keep shared platform code under `internal/apperr`, `internal/config`, and other small focused packages.
- Keep HTTP transport in `internal/api`.
- For new product areas, prefer feature-owned packages such as `internal/billing`, `internal/library`, or `internal/profile`.
- Put feature business logic and repository interfaces in the feature package.
- Add infrastructure implementations only when the feature exists.
- Do not introduce catch-all packages such as `util`, `helpers`, `common`, or a global `repository` layer.

## Dependency Direction

- `cmd/*` may depend on any `internal/*` package.
- `internal/api` may depend on feature packages and shared platform packages.
- Feature packages may depend on shared platform packages.
- Shared platform packages must not depend on feature packages or `internal/api`.

## Endpoint Pattern

- Register endpoints in `internal/api/<feature>.go` or a small feature subpackage when a route group becomes large.
- Keep handler input/output types close to the handler.
- Keep transport concerns in the API layer: HTTP status codes, request decoding, response encoding, auth context extraction.
- Move non-trivial business rules into a feature package rather than growing handlers.
- Prefer constructor injection through explicit dependency structs over package globals.

## Errors

- Return `apperr.Error` for domain and validation failures that should map to stable HTTP problem responses.
- Use stable machine-readable error codes such as `users.not_found` or `config.invalid_port`.
- Avoid exposing raw internal errors to clients.
- Wrap infrastructure failures with `apperr.Wrap` when they cross package boundaries.

## Auth

- Use `internal/authn.Verifier` as the seam for future Clerk integration.
- Keep Clerk SDK usage out of `internal/api` handlers directly.
- Protected endpoints should verify the incoming bearer token through the verifier and work with an `authn.Principal`, not Clerk-specific types.

## Data Access

- When sqlc is introduced, generated query code should stay close to concrete store implementations.
- Feature packages should define the repository interfaces they need.
- sqlc-backed stores should satisfy those feature-owned interfaces.
- Do not add a shared database abstraction before a real use case requires it.

## Comments

- Add comments only when behavior is non-obvious or when a boundary exists for a future integration.
- Do not add comments that restate the code.

## Testing

- Favor package-local tests that verify observable behavior.
- For API tests, assert status codes, response bodies, headers, and OpenAPI registration where relevant.
- Keep config tests table-free unless tables improve clarity.
- Update `shared/openapi.json` via `mise run schema` whenever API contracts change.

## Workflow

- Run `mise run qa` before finishing backend changes.
- Run `mise run schema` after endpoint or schema changes.
- Keep README and this file aligned with the actual package layout and runtime behavior.
