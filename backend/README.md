# Tonaris Backend

## Requirements

- Go
- Mise

## Architecture

- `cmd/server`: process bootstrap, config loading, logging, and graceful shutdown
- `cmd/gen-openapi`: OpenAPI generator for the shared frontend contract
- `internal/api`: HTTP app composition, middleware, system routes, and problem responses
- `internal/authn`: auth verification seam for future Clerk integration
- `internal/apperr`: application error model and HTTP status mapping
- `internal/config`: runtime configuration loading and validation
- `internal/projectpath`: module-root lookup used by local tooling

This backend is intentionally small. Keep handlers thin, avoid generic abstraction layers, and add feature packages only when there is real feature logic to own.

## Environment configuration

The backend reads:

- `.env` during local development
- process environment variables in production

Current variables:

- `TONARIS_ENV`
  - allowed values: `development`, `production`
  - default: `development`
- `PORT`
  - default: `8698`
  - stays unprefixed so Railway can provide it directly
- `CORS_ALLOWED_ORIGINS`
  - comma-separated list, for example `http://localhost:3000,https://app.example.com`
  - default: `http://localhost:3000` in development
  - required in production

Behavior and precedence:

- In development, `.env` is loaded if it exists.
- In production, `.env` is skipped when `TONARIS_ENV=production`.
- Process environment variables override values from `.env`.
- Missing `.env` is fine.
- Invalid values fail startup immediately.

See `AGENTS.md` for implementation patterns and package rules.

## Run locally

```bash
cp .env.example .env
mise run dev
```

The server listens on `http://localhost:8698` by default.

## Endpoints

- Hello: `GET /hello`
- Health: `GET /healthz`
- OpenAPI: `GET /api-doc/openapi.json`
- Docs: `GET /scalar`

## Development tasks

```bash
mise run fmt
mise run fmt-check
mise run lint
mise run test
mise run test-race
mise run cover
mise run build
mise run qa
```

`mise run qa` is the non-mutating validation pass. `mise run schema` is separate because it rewrites the generated OpenAPI file.

## Generate schema

```bash
mise run schema
```

This rewrites `../shared/openapi.json`.

## Future integrations

- Authentication should be added behind `internal/authn.Verifier`, with Clerk as the concrete provider later.
- Database access should use feature-owned interfaces with sqlc-backed implementations when persistence is added.
