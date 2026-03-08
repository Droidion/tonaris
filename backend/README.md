# Tonaris Backend

## Requirements

- Go
- Mise

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

Behavior and precedence:

- In development, `.env` is loaded if it exists.
- In production, `.env` is skipped when `TONARIS_ENV=production`.
- Process environment variables override values from `.env`.
- Missing `.env` is fine.
- Invalid values fail startup immediately.

## Run locally

```bash
cp .env.example .env
go run ./cmd/server
```

## Generate schema

```bash
go run ./cmd/gen-openapi
```

This rewrites `../shared/openapi.json`.

## Test

```bash
go test ./...
```
