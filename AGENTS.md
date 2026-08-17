# AGENTS.md

<!-- cspell:words fieldalignment golangci GOPATH languagetool purego -->
<!-- cspell:words shellcheck -->

This file provides guidance to AI agents working with code in this
repository. For an overview of what Shamash is and does, read
[README.md][readme] first.

## Repository Overview

Shamash is the e-mail component of the kagal.dev ecosystem, alongside
[kagal][kagal] itself. Transport is delegated to commercial providers;
mailboxes, filtering and search stay on infrastructure the operator
runs.

[README.md][readme] explains the design. What follows is only the part
of it that constrains how code gets written here:

- **Dovecot answers two questions, and only two.** How mail is
  represented on disc, and what the full-text search API looks like.
  A decision falling under either is settled by what Dovecot does, and
  a divergence is a defect even where the alternative looks defensible
  in isolation. Shamash need not run Dovecot to be bound by it; outside
  those two questions Dovecot carries no authority here.
- **The Solr endpoint follows from the second.** Dovecot's full-text
  search speaks to Solr, so that is the shape of the API — dictated from
  outside, not chosen. bleve, HNSW and bbolt are implementation detail
  behind it, held per user, and must not leak into the contract.
- **One binary, driven by cobra.** The mail server, the index owner and
  the command line are the same executable. New functionality is a
  cobra subcommand, not a second program.
- **Indexes are owned by a daemon; commands are RPC clients.** A user's
  search index admits one reader and writer at a time, so nothing opens
  one directly. The binary self-daemonises on demand, serves the index
  over RPC, and exits once idle. Any code path that reads or writes an
  index goes through that daemon — no exceptions, no direct handle as
  an optimisation.
- **Provider plugins are compiled in for now.** Drop-in loading is to
  follow, bound at run time through purego rather than the standard
  library's `plugin` package. `darvaza.org/x/tls/internal/macos` is the
  pattern to imitate, build-tag split included: portable declarations in
  one file, the platform implementation in another.

## Related Documentation

Shamash has no build documentation of its own — it uses the darvaza.org
shared build system verbatim, so the upstream references apply:

- [BUILDING.md][building] — build targets, tooling, linting, CI, and
  troubleshooting.
- [TESTING.md][testing] — testing patterns for darvaza.org projects.
- [README-coverage.md][coverage] — coverage system, in-tree.

## Prerequisites

- Go 1.25 or later — the module's declared minimum (`go.mod`).
- `make`, `git`, and a configured `$GOPATH`.

`golangci-lint` and `revive` are pinned by the `Makefile` and fetched on
demand. The Markdown and spelling tools (markdownlint, cspell,
languagetool, shellcheck) are optional: when absent they degrade to
no-ops, so a green local `make tidy` does not guarantee a green CI.

## Repository Layout

```text
shamash/
├── Makefile                # Shared darvaza.org build orchestration
├── go.mod                  # Go module: kagal.dev/shamash
├── .github/workflows/      # Build, platforms, codecov, renovate, Claude
├── internal/build/         # Build scripts and tool configuration
└── .tmp/                   # Generated artefacts (gitignored)
```

## Go

The Go side follows darvaza.org conventions, and `darvaza.org/core` is the
reference implementation to imitate when a pattern is unclear.

```bash
make all          # full cycle: get, generate, tidy, build
make tidy         # format, lint, validate
make test         # run tests (no cache reuse)
make coverage     # tests with coverage
make race         # tests with race detection
```

Points that most often catch agents out:

- Revive enforces hard limits: 40 lines per function, 5 arguments, 3
  results, cognitive complexity 7, cyclomatic complexity 10. Restructure
  the code; never relax the threshold.
- `fieldalignment` is enforced by CI. Use the probe-file workflow in
  [BUILDING.md][building] — running `-fix` over the tree strips every
  comment.
- Test files use `package <pkg>_test` and exercise the public surface.
  Table tests with two or more rows use `core.TestCase` plus
  `core.RunTestCases`.
- Build before lint; stubs need to compile first.

## TypeScript and npm

Shamash will also publish npm packages. None exist yet. When they are
added, follow the layout and conventions of [kagal][kagal]: a pnpm
workspace under `packages/`, ES modules throughout, `workspace:^` for
internal dependencies, Vitest for tests, and the split `tsconfig.json` /
`tsconfig.tools.json` / `tsconfig.tests.json` arrangement.

## Conventions

- British English in code, comments, documentation, and commit messages.
- Acronyms keep their case: `HTTP`, `URL`, `ID`, `SMTP` — never `Http`.
- Factory functions take a `new` or `make` prefix, not `create`.
- Treat every warning as an error. Fix the cause, or annotate a verified
  false positive at the source.
- Never stage in bulk. Use `git add` for new files only, then pass
  explicit paths to `git commit -s -F .tmp/commit-<slug>.txt`.
- No AI-attribution lines in commits, pull requests, or documents.
- Scratch files belong in `.tmp/`, never in the ambient `/tmp`.

## Pre-commit Checklist

1. Run `make tidy` until it passes, fixing what it reports.
2. Verify tests pass with `make test`.
3. Update this file and `README.md` when the workflow, behaviour, or API
   changes.

[building]: https://github.com/darvaza-proxy/core/blob/main/BUILDING.md
[coverage]: internal/build/README-coverage.md
[kagal]: https://github.com/kagal-dev/kagal
[readme]: README.md
[testing]: https://github.com/darvaza-proxy/core/blob/main/TESTING.md
