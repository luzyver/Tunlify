# Release Guide

## Versioning

This project uses [Calendar Versioning (CalVer)](https://calver.org/) with format `YYYY.MINOR.PATCH`.

- `YYYY` — release year
- `MINOR` — incremented for new features
- `PATCH` — incremented for bug fixes

## How to Release

```bash
git tag -a 2026.X.0 -m "Release 2026.X.0"
git push origin --tags
```

Then create a GitHub Release from the tag.

## Example

```
2026.1.0 → initial release
2026.2.0 → added project management
2026.2.1 → bug fix
2026.3.0 → next feature release
```
