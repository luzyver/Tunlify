# Release Guide

## Versioning

This project uses [Semantic Versioning](https://semver.org/) with automatic version detection from commit messages.

## Commit Convention

| Prefix | Bump | Example |
|--------|------|---------|
| `fix:` | patch (0.0.X) | `fix: resolve login timeout` |
| `docs:` | patch (0.0.X) | `docs: update README` |
| `chore:` | patch (0.0.X) | `chore: cleanup unused files` |
| `refactor:` | patch (0.0.X) | `refactor: simplify auth logic` |
| `feat:` | minor (0.X.0) | `feat: add metrics dashboard` |
| `feat!:` | major (X.0.0) | `feat!: redesign API auth flow` |
| `BREAKING CHANGE` in body | major (X.0.0) | Any commit with breaking change note |

## How to Release

1. Ensure all commits on `master` follow the convention above
2. Go to the **Actions** tab on GitHub
3. Select the **Release** workflow
4. Click **Run workflow**
5. Version is automatically determined from commits since the last tag

## Example Flow

```
v0.1.0 (last tag)
  ↓
fix: resolve restart error        → patch
feat: add metrics filtering       → minor (wins)
chore: update deps                → patch
  ↓
Release → v0.2.0
```

The highest bump wins: major > minor > patch.
