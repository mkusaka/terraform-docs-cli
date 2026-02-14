---
name: tfdc
description: Export Terraform provider documentation from the public registry to local files using the tfdc CLI. Use when you need to fetch, read, or reference Terraform provider docs (resources, data sources, guides, etc.) for any provider and version. Triggers on requests like "get Terraform docs", "export provider documentation", "fetch AWS provider resources docs", or when working with Terraform infrastructure that needs provider documentation locally.
---

# tfdc

## Overview

tfdc fetches Terraform provider documentation from the public registry and writes it to local files. Only `provider export` is currently implemented.

## Quick Start

Export a single provider's docs:

```bash
tfdc provider export -name aws -version 6.31.0 -out-dir ./terraform-docs
```

Export all providers from a lockfile:

```bash
tfdc -chdir=./infra provider export -out-dir ./terraform-docs
```

Read an exported doc:

```bash
cat ./terraform-docs/terraform/hashicorp/aws/6.31.0/docs/resources/instance.md
```

## provider export

Two modes are available.

### Legacy mode (single provider)

```bash
tfdc provider export \
  -name aws \
  -version 6.31.0 \
  -out-dir ./docs
```

Required: `-name`, `-version`, `-out-dir`

### Lockfile mode (multi-provider)

```bash
tfdc -chdir=./infra provider export -out-dir ./docs
```

Required: `-chdir` (global flag), `-out-dir`

Detects `.terraform.lock.hcl` in the `-chdir` directory and exports all listed providers. Filter to one provider with `-name`:

```bash
tfdc -chdir=./infra provider export -name aws -out-dir ./docs
```

### Optional flags

| Flag | Default | Description |
|---|---|---|
| `-namespace` | `hashicorp` | Provider namespace |
| `-format` | `markdown` | Output format (`markdown` or `json`) |
| `-categories` | `all` | Categories to export (comma-separated) |
| `-path-template` | See below | Output path template |
| `-clean` | off | Remove previous export before writing |

### Output layout

Default template: `{out}/terraform/{namespace}/{provider}/{version}/docs/{category}/{slug}.{ext}`

Example output: `docs/terraform/hashicorp/aws/6.31.0/docs/resources/instance.md`

Manifest: `docs/terraform/hashicorp/aws/6.31.0/docs/_manifest.json`

### Categories

`-categories all` expands to: `resources`, `data-sources`, `ephemeral-resources`, `functions`, `guides`, `overview`, `actions`, `list-resources`

Export specific categories:

```bash
tfdc provider export -name aws -version 6.31.0 -out-dir ./docs -categories resources,data-sources
```

### Path template placeholders

`{out}`, `{namespace}`, `{provider}`, `{version}`, `{category}`, `{slug}`, `{doc_id}`, `{ext}`

## Global flags

| Flag | Default | Description |
|---|---|---|
| `-chdir` | none | Switch working directory; auto-detects `.terraform.lock.hcl` |
| `-timeout` | `10s` | HTTP timeout |
| `-retry` | `3` | Retry count |
| `-registry-url` | `https://registry.terraform.io` | Registry base URL |
| `-insecure` | off | Skip TLS verification |
| `-user-agent` | `tfdc/dev` | User-Agent header |
| `-debug` | off | Debug logs to stderr |
| `-cache-dir` | `~/.cache/tfdc` | Cache directory |
| `-cache-ttl` | `24h` | Cache TTL |
| `-no-cache` | off | Disable cache |

## Exit codes

| Code | Meaning |
|---|---|
| `0` | Success |
| `1` | Invalid arguments or validation error |
| `2` | Not found |
| `3` | Remote API error |
| `4` | File write or serialization error |
