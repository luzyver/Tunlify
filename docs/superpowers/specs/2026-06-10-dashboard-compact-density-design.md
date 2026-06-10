# Dashboard Compact-Density Redesign — Design

**Date:** 2026-06-10
**Scope:** All 10 dashboard pages (behind login)
**Direction:** Compact density (option A) — same look, structure, and behavior; ~25–30% tighter
**Confirmed via:** visual calibration (moderate target approved)

## Goal

Fit more on screen with less scrolling by reducing spacing and element sizing.
Keep the existing visual identity (bitcoin/orange theme, fonts, glow background,
page-enter animation) and all functionality unchanged. This is a spacing/sizing
pass only — no logic, no component behavior, no color changes.

## Approach

Centralized first, then a per-page sweep:

1. **Tune the shared classes in `frontend/src/assets/main.css`.** Most spacing
   lives in `.page`, `.stat-*`, `.card-*`, `.btn-*`, `.empty`, and the `.v-table`
   overrides. Editing these cascades to all 10 pages from one source of truth and
   keeps the result consistent.
2. **Per-page sweep** for spacing that does not cascade: Vuetify margin utilities
   (`mb-8`, `mb-6`, …) and a few bulky inline paddings in the view templates.

Rejected alternatives: per-page-only edits (inconsistent, more churn); separate
`--compact` class variants (two systems to maintain).

## Concrete changes — `main.css`

| Selector | Property | Now | Compact |
|---|---|---|---|
| `.page` | padding | 24px 32px | 16px 24px |
| `.page-header` | margin-bottom | 24px | 16px |
| `.page-title` | font-size | 26px | 22px |
| `.stat-grid` | gap | 16px | 12px |
| `.stat-grid` | margin-bottom | 24px | 16px |
| `.stat-grid` | min column (minmax) | 200px | 168px |
| `.stat-card` | padding | 20px 24px | 14px 16px |
| `.stat-label` | margin-bottom | 8px | 6px |
| `.stat-value` | font-size | 24px | 20px |
| `.card-header` | padding | 16px 24px | 12px 16px |
| `.card-body` | padding | 24px | 16px |
| `.empty` | padding | 32px | 20px |
| `.btn-lg` | height | 44px | 40px |
| `.v-table` th/td | vertical cell padding | default | tighter (≈6px) |

## Concrete changes — per page (template sweep)

Standardize section spacing to a tighter scale across all views:

- `mb-8` (32px) → `mb-5` (20px)
- `mb-6` (24px) → `mb-4` (16px)
- Leave `mb-4`/`mb-3`/`mb-2` as-is unless visibly loose.
- Trim the few bulky inline `padding:`/`margin:` blocks in view templates to match
  the compact card paddings.

`DataTable.vue`:
- toolbar wrapper `mb-3` → `mb-2`
- empty-state cell `py-8` → `py-6`
- (table already uses `density="compact"`; cell padding handled via the `.v-table`
  override above)

Pages in scope (all use the `.page` wrapper): Dashboard, Projects, Config/Ingress,
Health, Metrics, Backups, Audit, Notifications, TCP Access, Settings.

## Out of scope

- No color, font, or icon changes.
- No layout restructuring (no consolidating sections/pages, no two-column rework —
  that was option B, not chosen).
- No new components or abstractions.
- Background glows and the `pg-enter` animation stay.

## Verification

- `npm run build` succeeds (no template/type errors).
- Manual spot check: every page still renders with the same content and controls;
  visibly tighter; no clipped/overlapping elements; tables and pagination intact.
- No functional regressions (purely presentational diff).
