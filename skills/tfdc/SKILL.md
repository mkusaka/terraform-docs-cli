---
name: tfdc
description: Retrieve Terraform documentation from the public registry using the tfdc CLI. Supports provider doc search/get/export, module search/get, policy search/get, and Terraform style/module-dev guides. Use when you need to find, read, or export Terraform provider docs, module docs, policy docs, or development guides. Triggers on requests like "get Terraform docs", "search AWS provider resources", "fetch module details", "look up policy", or "show Terraform style guide".
---

# tfdc

## Overview

tfdc retrieves Terraform documentation from the public registry. It supports provider doc search/get/export, module search/get, policy search/get, and Terraform guides.

All search/get/guide commands support `-format text|json|markdown` (default: `text`).

## Workflow

Provider docs use a two-step flow: search for candidate `provider_doc_id` values, then get full content by ID.

```bash
# Step 1: Find doc IDs
tfdc provider search -name aws -service ec2 -type resources -format json

# Step 2: Fetch full content
tfdc provider get -doc-id 10595066
```

Module and policy docs follow the same search-then-get pattern.

## provider search

Search provider documentation by service slug.

```bash
tfdc provider search \
  -name aws \
  -service ec2 \
  -type resources \
  [-namespace hashicorp] \
  [-version latest] \
  [-limit 20] \
  [-format text]
```

Required: `-name`, `-service`, `-type`

`-type` values: `resources`, `data-sources`, `ephemeral-resources`, `functions`, `guides`, `overview`, `actions`, `list-resources`

Output fields: `provider_doc_id`, `title`, `category`, `description`, `provider`, `namespace`, `version`

## provider get

Fetch full provider doc content by `provider_doc_id`.

```bash
tfdc provider get -doc-id 10595066 [-format text]
```

Required: `-doc-id` (numeric)

JSON output: `{ "id", "content", "content_type" }`

## provider export

Export all docs for a provider version to local files.

```bash
# Single provider
tfdc provider export -name aws -version 6.31.0 -out-dir ./docs

# From lockfile
tfdc -chdir=./infra provider export -out-dir ./docs
```

Required: `-name` + `-version` + `-out-dir` (legacy mode) or `-chdir` + `-out-dir` (lockfile mode)

| Flag | Default | Description |
|---|---|---|
| `-namespace` | `hashicorp` | Provider namespace |
| `-format` | `markdown` | Persist format (`markdown` or `json`) |
| `-categories` | `all` | Categories to export (comma-separated) |
| `-path-template` | See below | Output path template |
| `-clean` | off | Remove previous export before writing |

Default path: `{out}/terraform/{namespace}/{provider}/{version}/docs/{category}/{slug}.{ext}`

Manifest: `{out}/terraform/{namespace}/{provider}/{version}/docs/_manifest.json`

## module search

Search the Terraform module registry.

```bash
tfdc module search -query vpc [-offset 0] [-limit 20] [-format text]
```

Required: `-query`

Output fields: `module_id`, `name`, `description`, `downloads`, `verified`, `published_at`

## module get

Fetch module details by exact module ID.

```bash
tfdc module get -id terraform-aws-modules/vpc/aws/6.0.1 [-format text]
```

Required: `-id` (format: `namespace/name/provider/version`)

JSON output: `{ "id", "content", "content_type" }`

## policy search

Search Terraform policy sets.

```bash
tfdc policy search -query cis [-format text]
```

Required: `-query`

Output fields: `terraform_policy_id`, `name`, `title`, `downloads`

## policy get

Fetch policy details by exact policy ID.

```bash
tfdc policy get -id policies/hashicorp/CIS-Policy-Set-for-AWS-Terraform/1.0.1 [-format text]
```

Required: `-id` (format: `policies/namespace/name/version`)

JSON output: `{ "id", "content", "content_type" }`

## guide style

Fetch the Terraform style guide.

```bash
tfdc guide style [-format text]
```

## guide module-dev

Fetch the Terraform module development guide.

```bash
tfdc guide module-dev [-section all] [-format text]
```

`-section` values: `all`, `index`, `composition`, `structure`, `providers`, `publish`, `refactoring`

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

## JSON output contract

Search commands:

```json
{
  "items": [...],
  "total": 0
}
```

Detail/get/guide commands:

```json
{
  "id": "string",
  "content": "string",
  "content_type": "text/markdown"
}
```
