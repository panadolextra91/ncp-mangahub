---
phase: 12
plan: 3
wave: 3
autonomous: true
files_modified:
  - cmd/mangahub/main.go
  - internal/interfaces/cli/root.go
  - internal/interfaces/cli/auth.go
  - internal/interfaces/cli/manga.go
  - internal/interfaces/cli/library.go
---

# Plan: Cobra CLI Implementation

## Tasks

<task>
<action>
Initialize the Cobra CLI structure in `cmd/mangahub`.
- Install `github.com/spf13/cobra`.
- Create `root.go` with basic configuration (Server URL, Config file handling).
</action>
<acceptance_criteria>
- `cmd/mangahub/main.go` exists.
- `go run cmd/mangahub/main.go --help` prints the help message.
</acceptance_criteria>
</task>

<task>
<read_first>
- `docs/cli_manual.md`
</read_first>
<action>
Implement the `auth` subcommand group.
- `mangahub auth login`
- `mangahub auth register`
- Save JWT to `~/.mangahub/config.json`.
</action>
<acceptance_criteria>
- `mangahub auth login` successfully authenticates and stores the token.
</acceptance_criteria>
</task>

<task>
<read_first>
- `docs/cli_manual.md`
</read_first>
<action>
Implement the `manga` subcommand group.
- `mangahub manga search <query>`
- `mangahub manga info <id>`
- Print results in a clean table format.
</action>
<acceptance_criteria>
- `mangahub manga search "One Piece"` returns a table of matching manga.
</acceptance_criteria>
</task>

<task>
<read_first>
- `docs/cli_manual.md`
</read_first>
<action>
Implement the `library` subcommand group.
- `mangahub library add <id> --status reading`
- `mangahub library list`
- `mangahub library update <id> --chapter <n>`
</action>
<acceptance_criteria>
- `mangahub library add` updates the user's progress in the DB.
- `mangahub library list` shows the user's reading list.
</acceptance_criteria>
</task>

<task>
<action>
Implement CLI "Golden Path" integration tests.
- **Scenario:** Register -> Login -> Search -> Add to Library -> Update Progress.
- **Edge Case:** Command run without server active (should show "Server unreachable" not a panic).
- **Edge Case:** Invalid credentials in `auth login`.
- **Edge Case:** Invalid Manga ID in `library add`.
</action>
<acceptance_criteria>
- A test script or Go test file verifies the full CLI lifecycle.
- Error messages match the `cli_manual.md` requirements.
</acceptance_criteria>
</task>
