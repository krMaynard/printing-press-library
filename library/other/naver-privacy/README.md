# Naver Privacy CLI

**Naver's transparency-report series — Korea's first, published since 2015 — as tidy per-category statistics instead of a CMS page blob.**

The Naver Privacy Center ships its entire 2012-present government data-request series (warrants, communication data, restriction measures, log records) inside one CMS page payload. This CLI unpacks it: `statistics` emits tidy per-period, per-category rows, `summaries` extracts Naver's own commentary, and `whitepapers` indexes the downloadable NAVER Privacy Report library.

Learn more at [Naver Privacy](https://privacy.naver.com).

Created by [@krMaynard](https://github.com/krMaynard) (Kieran Maynard).

## Install

The recommended path installs both the `naver-privacy-pp-cli` binary and the `pp-naver-privacy` agent skill (Claude Code, Codex, Cursor, Gemini CLI, GitHub Copilot, and other agents supported by the upstream [`skills`](https://github.com/vercel-labs/skills) CLI) in one shot:

```bash
npx -y @mvanhorn/printing-press-library install naver-privacy
```

For CLI only (no skill):

```bash
npx -y @mvanhorn/printing-press-library install naver-privacy --cli-only
```

For skill only — installs the skill into the same agents as the default command above, but skips the CLI binary (use this to update or reinstall just the skill):

```bash
npx -y @mvanhorn/printing-press-library install naver-privacy --skill-only
```

To constrain the skill install to one or more specific agents (repeatable — agent names match the [`skills`](https://github.com/vercel-labs/skills) CLI):

```bash
npx -y @mvanhorn/printing-press-library install naver-privacy --agent claude-code
npx -y @mvanhorn/printing-press-library install naver-privacy --agent claude-code --agent codex
```

### Without Node (Go fallback)

If `npx` isn't available (no Node, offline), install the CLI directly via Go (requires Go 1.26.4 or newer):

```bash
go install github.com/mvanhorn/printing-press-library/library/other/naver-privacy/cmd/naver-privacy-pp-cli@latest
```

This installs the CLI only — no skill.

### Pre-built binary

Download a pre-built binary for your platform from the [latest release](https://github.com/mvanhorn/printing-press-library/releases/tag/naver-privacy-current). On macOS, clear the Gatekeeper quarantine: `xattr -d com.apple.quarantine <binary>`. On Unix, mark it executable: `chmod +x <binary>`.

<!-- pp-hermes-install-anchor -->
## Install for Hermes

Install the CLI binary first. The installer writes binaries to a per-user managed bin directory by default: `$HOME/.local/bin` on macOS/Linux and `%LOCALAPPDATA%\Programs\PrintingPress\bin` on Windows.

```bash
npx -y @mvanhorn/printing-press-library install naver-privacy --cli-only
```

Then install the focused Hermes skill.

From the Hermes CLI:

```bash
hermes skills install mvanhorn/printing-press-library/cli-skills/pp-naver-privacy --force
```

Inside a Hermes chat session:

```bash
/skills install mvanhorn/printing-press-library/cli-skills/pp-naver-privacy --force
```

Restart the Hermes session or gateway if the newly installed skill is not visible immediately.

## Install for OpenClaw
Install both the CLI binary and the focused OpenClaw skill. The installer defaults binaries to a per-user bin directory (`$HOME/.local/bin` on macOS/Linux, `%LOCALAPPDATA%\Programs\PrintingPress\bin` on Windows):

```bash
npx -y @mvanhorn/printing-press-library install naver-privacy --agent openclaw
```

Restart the OpenClaw session or gateway if the newly installed skill is not visible immediately.

## Use with Claude Desktop

This CLI ships an [MCPB](https://github.com/modelcontextprotocol/mcpb) bundle — Claude Desktop's standard format for one-click MCP extension installs (no JSON config required).

To install:

1. Download the `.mcpb` for your platform from the [latest release](https://github.com/mvanhorn/printing-press-library/releases/tag/naver-privacy-current).
2. Double-click the `.mcpb` file. Claude Desktop opens and walks you through the install.

Requires Claude Desktop 1.0.0 or later. Pre-built bundles ship for macOS Apple Silicon (`darwin-arm64`) and Windows (`amd64`, `arm64`); for other platforms, use the manual config below.

<details>
<summary>Manual JSON config (advanced)</summary>

If you can't use the MCPB bundle (older Claude Desktop, unsupported platform), install the MCP binary and configure it manually.


```bash
go install github.com/mvanhorn/printing-press-library/library/other/naver-privacy/cmd/naver-privacy-pp-mcp@latest
```

Add to your Claude Desktop config (`~/Library/Application Support/Claude/claude_desktop_config.json`):

```json
{
  "mcpServers": {
    "naver-privacy": {
      "command": "naver-privacy-pp-mcp"
    }
  }
}
```

</details>

## Quick Start

```bash
# Verify connectivity — the API is public, no credentials needed
naver-privacy-pp-cli doctor --dry-run

# The full transparency series as tidy CSV rows
naver-privacy-pp-cli statistics --csv

# The raw CMS payload if you need every field
naver-privacy-pp-cli pages get TRANSPARENCY_REPORT_STATISTICS

# Current Privacy Center notices
naver-privacy-pp-cli notices list

```

## Unique Features

These capabilities aren't available in any other tool for this API.

### Whole-archive analytics
- **`statistics`** — The full 2012-present transparency series as one tidy table — one row per half-year and request category with request/processed/provided/average counts and compliance rates.

  _Reach for this for any question about Naver's government data-request trends — it replaces hand-parsing specificAreaJson._

  ```bash
  naver-privacy-pp-cli statistics --category warrant --csv
  ```
- **`whitepapers`** — List the NAVER Privacy Report and privacy-whitepaper publications (titles, authors, attached documents) from the report library.

  _Use this to find Naver's own published analyses and the downloadable report documents behind the statistics._

  ```bash
  naver-privacy-pp-cli whitepapers --agent
  ```
- **`summaries`** — Naver's own per-half-year commentary on the statistics, HTML-stripped to plain text.

  _Use this before interpreting a spike or drop in the numbers — Naver usually explains it in the period's own commentary._

  ```bash
  naver-privacy-pp-cli summaries --year 2025
  ```

## Recipes


### Warrant series for charting

```bash
naver-privacy-pp-cli statistics --category warrant --agent --select year,period,requests,processed,rate
```

Tidy JSON of search-and-seizure warrant volumes and compliance rates per half-year since 2012, narrowed to chartable fields.

### Why did this period spike?

```bash
naver-privacy-pp-cli summaries --year 2018
```

Naver's own plain-text commentary for 2018's half-years (e.g. the Camp Mobile merger's effect on the numbers).

### Index the report library

```bash
naver-privacy-pp-cli whitepapers --agent
```

Every NAVER Privacy Report / whitepaper publication with its attached document paths.

### Raw CMS payload

```bash
naver-privacy-pp-cli pages get TRANSPARENCY_REPORT_STATISTICS --agent --select specificAreaJson.statistics.year,specificAreaJson.statistics.period,specificAreaJson.statistics.warrantRequestCount
```

Drop to the raw page payload and narrow it with dotted --select paths when a field is not in the tidy view.

## Usage

Run `naver-privacy-pp-cli --help` for the full command reference and flag list.

## Commands

### notices

Privacy Center notices.

- **`naver-privacy-pp-cli notices get`** - Returns one notice with its full HTML body and prev/next navigation stubs.
- **`naver-privacy-pp-cli notices list`** - Returns notices, newest first. Paginated with 0-indexed page and size.

### pages

CMS pages and section intros — including the transparency-report statistics series carried by the TRANSPARENCY_REPORT_STATISTICS page.

- **`naver-privacy-pp-cli pages get`** - Returns a CMS page's rendered HTML plus its structured
`specificAreaJson` payload. For `TRANSPARENCY_REPORT_STATISTICS` the
payload's `statistics` array is the full transparency-report series:
one entry per half-year since 2012, each carrying request, processing,
provided-item, and per-request-average counts plus compliance rates
for search & seizure warrants, communication user information,
communication-restricting measures, and communication confirmation
data.
- **`naver-privacy-pp-cli pages get-intro`** - Returns the intro blurb (title, text, hero image) for a Privacy Center section. `REPORT` is the transparency-report section; `ACTIVITY` is the protection-activities section.


## Output Formats

```bash
# Human-readable table (default in terminal, JSON when piped)
naver-privacy-pp-cli notices list

# JSON for scripting and agents
naver-privacy-pp-cli notices list --json

# Filter to specific fields
naver-privacy-pp-cli notices list --json --select id,name,status

# Dry run — show the request without sending
naver-privacy-pp-cli notices list --dry-run

# Agent mode — JSON + compact + no prompts in one flag
naver-privacy-pp-cli notices list --agent
```

## Agent Usage

This CLI is designed for AI agent consumption:

- **Non-interactive** - never prompts, every input is a flag
- **Pipeable** - `--json` output to stdout, errors to stderr
- **Filterable** - `--select id,name` returns only fields you need
- **Previewable** - `--dry-run` shows the request without sending
- **Read-only by default** - this CLI does not create, update, delete, publish, send, or mutate remote resources
- **Offline-friendly** - sync/search commands can use the local SQLite store when available
- **Agent-safe by default** - no colors or formatting unless `--human-friendly` is set

Exit codes: `0` success, `2` usage error, `3` not found, `5` API error, `7` rate limited, `10` config error.

## Health Check

```bash
naver-privacy-pp-cli doctor
```

Verifies configuration and connectivity to the API.

## Configuration

Config file: `~/.config/naver-privacy-pp-cli/config.toml`

Static request headers can be configured under `headers`; per-command header overrides take precedence.

## Troubleshooting
**Not found errors (exit code 3)**
- Check the resource ID is correct
- Run the `list` command to see available items

### API-specific
- **pages get returns 400 BAD_REQUEST** — The page code is unknown to the CMS — codes are case-sensitive strings like TRANSPARENCY_REPORT_STATISTICS or NAVER_REPORT; check the code against the Privacy Center site.
- **Statistics cells contain '-'** — '-' marks categories Naver no longer complies with (no communication user information provided since October 2012); treat it as not-applicable, not zero.
