# Codebase Cleanup Refactor — Design

**Date:** 2026-06-10
**Scope:** Backend (Go) + Frontend (Vue) + repo hygiene
**Depth:** Conservative — **no runtime behavior change**

## Goal

Reduce duplication and noise without altering behavior. Every change should read
as "same thing, less code." Reviewable as a mechanical clean-up.

## Verification

- Backend: `go build ./...`, `go vet ./...`, `gofmt -l .` (empty).
- Frontend: `npm run build` (esbuild via Vite; no separate type-check script exists).

## Work items

### 1. Repo hygiene (first commit)
- Add `.gitattributes` with `* text=auto eol=lf` plus binary exceptions, so the
  CRLF↔LF churn (all 23 currently-"modified" files are pure line-ending flips)
  stops permanently.
- Renormalize tracked files to LF in one dedicated commit.
- Run `gofmt -w` on backend (the 14 `gofmt -l` hits are CRLF artifacts).
- Frontend has no formatter/linter configured; do **not** introduce one (would be
  a large diff + new dep — out of scope for conservative).

### 2. Backend response helpers (main dedup)
- New `internal/handler/respond.go`:
  - `writeJSON(w http.ResponseWriter, status int, v any)`
  - `writeError(w http.ResponseWriter, status int, msg string)`
- Replace the 34 manual `Content-Type` + `json.NewEncoder(w).Encode(...)`
  sequences and 22 hardcoded `` `{"error":...}` `` strings across handlers.
- Preserve existing status codes and JSON shapes exactly.

### 3. Backend: hand-rolled reader + ignored errors
- `notifications.go`: delete `jsonReader`/`jsonBody` custom `io.Reader`; use
  `bytes.NewReader`. Behavior-preserving (and removes a fragile EOF impl).
- Surface only the silently-ignored errors that can mask real failures
  (`db.Exec` in `NewAuth`, `os.WriteFile`/`os.MkdirAll` in notifications) via
  `log.Printf`. Harmless ignores (context type-asserts, token gen) stay as-is to
  keep the diff tight. No control-flow change.

### 4. Frontend: error surfacing
- Replace empty `catch {}` blocks (which silently hide failures) with consistent
  surfacing via the existing `actionLog` store. No new abstractions, no view
  restructuring.

## Explicitly out of scope (conservative)
- No package reorg, no splitting large views (LoginView/ProjectsView).
- No new layers/interfaces, no dependency changes.
- No generic frontend `useResource` composable (fetch pattern too light/varied
  to justify without a moderate restructure).
