# Utilities Catalogue

This repo gathers independently built utilities and publishes an index at https://utils.rselbach.com.

## Register a Utility
- create a top-level directory (for example `util-one/`)
- add `util-one/util.yaml` with:

```yaml
name: My Utility
description: what the utility provides
slug: util-one # optional; defaults to the directory name
```

## Local Catalogue Preview

```bash
go run ./cmd/catalog -base-url https://utils.rselbach.com -out site/index.html
```

Open `site/index.html` to review the rendered list.

## Automation
- GitHub Actions workflow `.github/workflows/catalogue.yml` runs on `main` pushes
- the workflow executes `go test`, regenerates `site/index.html`, and deploys to GitHub Pages
- override the published base URL by setting the repo variable `UTILS_BASE_URL`
