# Utils Collection

Automated hosting for multiple independent utilities served from a single domain.

## Structure

```
/
├── .github/workflows/
│   └── deploy.yml          # CI/CD workflow
├── scripts/
│   ├── generate_index.py   # Generates main index page
│   └── build_utility.py    # Builds individual utilities
├── index-template.html     # Template for main page
├── example-util/
│   ├── util.json          # Metadata
│   └── index.html         # Utility content
└── README.md
```

## Adding a New Utility

1. Create directory: `mkdir my-utility`
2. Add `util.json`:
```json
{
  "name": "My Utility",
  "description": "What this utility does",
  "build": "npm run build"
}
```
3. Add your utility files (static HTML or source that builds to HTML)
4. Commit and push

The `build` field is optional. Omit it for static utilities.

## Metadata Schema

### Required Fields
- `name`: Display name
- `description`: Brief description

### Optional Fields
- `build`: Shell command to run during CI/CD

## Local Development

### Generate Index Locally
```bash
python3 scripts/generate_index.py
```

### Build a Utility Locally
```bash
python3 scripts/build_utility.py <utility-dir>
```

## GitHub Secrets Setup

Configure these secrets in repository settings:

- `DEPLOY_SSH_KEY`: Private SSH key for server access
- `DEPLOY_HOST`: Server hostname or IP
- `DEPLOY_PATH`: Target directory (e.g., `/var/www/utils`)
- `DEPLOY_USER`: SSH username

### Generating SSH Key

```bash
ssh-keygen -t ed25519 -C "github-actions" -f deploy_key
```

Add `deploy_key.pub` to server's `~/.ssh/authorized_keys`, then add `deploy_key` contents to `DEPLOY_SSH_KEY` secret.

## Caddy Configuration

Manual Caddy config for utils.rselbach.com (example):

```caddy
utils.rselbach.com {
    root * /var/www/utils
    file_server
    try_files {path} {path}/index.html {path}.html
}
```

Add to `/etc/caddy/conf.d/utils.conf` and reload: `sudo systemctl reload caddy`

## CI/CD Workflow

On push to main:
1. Detects changed utility directories
2. Builds only changed utilities
3. Generates index.html with all utilities
4. Deploys via rsync to server

Manual trigger: Builds all utilities regardless of changes

## Examples

### Static Utility
```
my-static-util/
├── util.json
├── index.html
├── style.css
└── script.js
```

util.json:
```json
{
  "name": "My Static Tool",
  "description": "Simple static page"
}
```

### Built Utility (Node.js)
```
my-built-util/
├── util.json
├── package.json
├── src/
└── dist/
```

util.json:
```json
{
  "name": "My Built Tool",
  "description": "Requires build step",
  "build": "npm ci && npm run build"
}
```

Build output should go to a location served by web server (typically project root or `dist/` with appropriate routing).
