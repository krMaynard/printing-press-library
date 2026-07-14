---
name: pp-naver-privacy
description: "Naver's transparency-report series — Korea's first, published since 2015 — as tidy per-category statistics instead of a CMS page blob. Trigger phrases: `naver transparency report`, `korean government data requests naver`, `naver warrant statistics`, `naver privacy center`, `use naver-privacy`, `run naver-privacy`."
author: "Kieran Maynard"
license: "Apache-2.0"
argument-hint: "<command> [args] | install cli|mcp"
allowed-tools: "Read Bash"
regions: ["KR"]
api_language: "ko"
metadata:
  openclaw:
    requires:
      bins:
        - naver-privacy-pp-cli
    install:
      - kind: go
        bins: [naver-privacy-pp-cli]
        module: github.com/mvanhorn/printing-press-library/library/other/naver-privacy/cmd/naver-privacy-pp-cli
---

# Naver Privacy — Printing Press CLI

## Prerequisites: Install the CLI

This skill drives the `naver-privacy-pp-cli` binary. **You must verify the CLI is installed before invoking any command from this skill.** If it is missing, install it first:

1. Install via the Printing Press installer. It defaults binaries to `$HOME/.local/bin` on macOS/Linux and `%LOCALAPPDATA%\Programs\PrintingPress\bin` on Windows:
   ```bash
   npx -y @mvanhorn/printing-press-library install naver-privacy --cli-only
   ```
2. Verify: `naver-privacy-pp-cli --version`
3. Ensure the reported install directory is on `$PATH` for the agent/runtime that will invoke this skill.

If the `npx` install fails (no Node, offline, etc.), fall back to a direct Go install (requires Go 1.26.4 or newer). This installs into `$GOPATH/bin` (default `$HOME/go/bin`), so add that directory to `$PATH` instead:

```bash
go install github.com/mvanhorn/printing-press-library/library/other/naver-privacy/cmd/naver-privacy-pp-cli@latest
```

If `--version` reports "command not found" after install, the runtime cannot see the binary directory on `$PATH`. Do not proceed with skill commands until verification succeeds.

The Naver Privacy Center ships its entire 2012-present government data-request series (warrants, communication data, restriction measures, log records) inside one CMS page payload. This CLI unpacks it: `statistics` emits tidy per-period, per-category rows, `summaries` extracts Naver's own commentary, and `whitepapers` indexes the downloadable NAVER Privacy Report library.

## When to Use This CLI

Use this CLI for questions about South Korean government requests for Naver user data — warrant volumes, compliance rates, log-record requests — and for Naver's own privacy publications. It is the only machine-friendly path to the full statistics series.

## Anti-triggers

Do not use this CLI for:
- Do not use it for Naver's developer APIs (search, maps, Papago, CLOVA) — this wraps only the Privacy Center content API
- Do not use it for content-moderation or copyright-takedown statistics; Naver's transparency report covers government data requests only
- Content is Korean-only — do not expect localized English field values

## Unique Capabilities

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

## Command Reference

**notices** — Privacy Center notices.

- `naver-privacy-pp-cli notices get` — Returns one notice with its full HTML body and prev/next navigation stubs.
- `naver-privacy-pp-cli notices list` — Returns notices, newest first. Paginated with 0-indexed page and size.

**pages** — CMS pages and section intros — including the transparency-report statistics series carried by the TRANSPARENCY_REPORT_STATISTICS page.

- `naver-privacy-pp-cli pages get` — Returns a CMS page's rendered HTML plus its structured `specificAreaJson` payload.
- `naver-privacy-pp-cli pages get-intro` — Returns the intro blurb (title, text, hero image) for a Privacy Center section.


### Finding the right command

When you know what you want to do but not which command does it, ask the CLI directly:

```bash
naver-privacy-pp-cli which "<capability in your own words>"
```

`which` resolves a natural-language capability query to the best matching command from this CLI's curated feature index. Exit code `0` means at least one match; exit code `2` means no confident match — fall back to `--help` or use a narrower query.

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

## Auth Setup

No authentication required.

Run `naver-privacy-pp-cli doctor` to verify setup.

## Agent Mode

Add `--agent` to any command. Expands to: `--json --compact --no-input --no-color --yes`.

- **Pipeable** — JSON on stdout, errors on stderr
- **Filterable** — `--select` keeps a subset of fields. Dotted paths descend into nested structures; arrays traverse element-wise. Critical for keeping context small on verbose APIs:

  ```bash
  naver-privacy-pp-cli notices list --agent --select id,name,status
  ```
- **Previewable** — `--dry-run` shows the request without sending
- **Offline-friendly** — sync/search commands can use the local SQLite store when available
- **Non-interactive** — never prompts, every input is a flag
- **Read-only** — do not use this CLI for create, update, delete, publish, comment, upvote, invite, order, send, or other mutating requests

### Response envelope

Commands that read from the local store or the API wrap output in a provenance envelope:

```json
{
  "meta": {"source": "live" | "local", "synced_at": "...", "reason": "..."},
  "results": <data>
}
```

Parse `.results` for data and `.meta.source` to know whether it's live or local. A human-readable `N results (live)` summary is printed to stderr only when stdout is a terminal AND no machine-format flag (`--json`, `--csv`, `--compact`, `--quiet`, `--plain`, `--select`) is set — piped/agent consumers and explicit-format runs get pure JSON on stdout.

## Agent Feedback

When you (or the agent) notice something off about this CLI, record it:

```
naver-privacy-pp-cli feedback "the --since flag is inclusive but docs say exclusive"
naver-privacy-pp-cli feedback --stdin < notes.txt
naver-privacy-pp-cli feedback list --json --limit 10
```

Entries are stored locally at `~/.local/share/naver-privacy-pp-cli/feedback.jsonl`. They are never POSTed unless `NAVER_PRIVACY_FEEDBACK_ENDPOINT` is set AND either `--send` is passed or `NAVER_PRIVACY_FEEDBACK_AUTO_SEND=true`. Default behavior is local-only.

Write what *surprised* you, not a bug report. Short, specific, one line: that is the part that compounds.

## Output Delivery

Every command accepts `--deliver <sink>`. The output goes to the named sink in addition to (or instead of) stdout, so agents can route command results without hand-piping. Three sinks are supported:

| Sink | Effect |
|------|--------|
| `stdout` | Default; write to stdout only |
| `file:<path>` | Atomically write output to `<path>` (tmp + rename) |
| `webhook:<url>` | POST the output body to the URL (`application/json` or `application/x-ndjson` when `--compact`) |

Unknown schemes are refused with a structured error naming the supported set. Webhook failures return non-zero and log the URL + HTTP status on stderr.

## Named Profiles

A profile is a saved set of flag values, reused across invocations. Use it when a scheduled agent calls the same command every run with the same configuration - HeyGen's "Beacon" pattern.

```
naver-privacy-pp-cli profile save briefing --json
naver-privacy-pp-cli --profile briefing notices list
naver-privacy-pp-cli profile list --json
naver-privacy-pp-cli profile show briefing
naver-privacy-pp-cli profile delete briefing --yes
```

Explicit flags always win over profile values; profile values win over defaults. `agent-context` lists all available profiles under `available_profiles` so introspecting agents discover them at runtime.

## Exit Codes

| Code | Meaning |
|------|---------|
| 0 | Success |
| 2 | Usage error (wrong arguments) |
| 3 | Resource not found |
| 5 | API error (upstream issue) |
| 7 | Rate limited (wait and retry) |
| 10 | Config error |

## Argument Parsing

Parse `$ARGUMENTS`:

1. **Empty, `help`, or `--help`** → show `naver-privacy-pp-cli --help` output
2. **Starts with `install`** → ends with `mcp` → MCP installation; otherwise → see Prerequisites above
3. **Anything else** → Direct Use (execute as CLI command with `--agent`)

## MCP Server Installation

1. Install the MCP server:
   ```bash
   go install github.com/mvanhorn/printing-press-library/library/other/naver-privacy/cmd/naver-privacy-pp-mcp@latest
   ```
2. Register with Claude Code:
   ```bash
   claude mcp add naver-privacy-pp-mcp -- naver-privacy-pp-mcp
   ```
3. Verify: `claude mcp list`

## Direct Use

1. Check if installed: `which naver-privacy-pp-cli`
   If not found, offer to install (see Prerequisites at the top of this skill).
2. Match the user query to the best command from the Unique Capabilities and Command Reference above.
3. Execute with the `--agent` flag:
   ```bash
   naver-privacy-pp-cli <command> [subcommand] [args] --agent
   ```
4. If ambiguous, drill into subcommand help: `naver-privacy-pp-cli <command> --help`.
