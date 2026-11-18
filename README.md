# StatBoard

A fast, read-only command-line dashboard for the status of all your Git repositories at once.

If you juggle several projects, getting back up to speed means running `git status` in each one. StatBoard replaces that with a single command that shows every tracked repo's branch, uncommitted changes, and sync state in one table. It never modifies a repository.

<!-- Add a screenshot of the table here, e.g. ![StatBoard screenshot](docs/screenshot.png) -->

## Features

- **One table for every repo:** name, path, branch, working-tree changes, remote sync, and a color-coded status
- **Detailed change counts:** staged, modified, and untracked files are shown separately instead of a generic "dirty" flag
- **Ahead/behind tracking:** how many commits each branch is ahead of or behind its remote-tracking branch
- **Filters:** `--dirty` for repos with uncommitted changes, `--behind` for repos behind their remote
- **Honest about unknowns:** a repo with no upstream, or one that can't be checked, is shown as `UNKNOWN` or `ERROR`, never as "up to date"
- **Resilient:** one broken or moved repo shows an error row while the rest still display
- **Fast:** repos are checked concurrently (up to 8 at a time)
- **Read-only:** StatBoard never fetches, pulls, commits, or changes files

## Requirements

- [Go](https://go.dev/dl/) (see `go.mod` for the minimum version)
- [Git](https://git-scm.com/downloads) available on your `PATH`

## Setup

### Option 1: Install with Go

```bash
go install github.com/<your-username>/statboard@latest
```

### Option 2: Build from source

```bash
git clone https://github.com/<your-username>/statboard.git
cd statboard
go install .
```

`go install` places the `statboard` binary in `$(go env GOPATH)/bin`. Make sure that folder is on your `PATH` so you can run `statboard` from anywhere.

To try it without installing, run `go run .` from the project folder instead.

## Usage

Add the repositories you want to track (paths must be Git repos):

```bash
statboard config add ~/projects/my-app
statboard config add ~/projects/another-repo
statboard config list
statboard config remove ~/projects/another-repo
```

Show the dashboard:

```bash
statboard                  # all tracked repos
statboard --dirty          # only repos with uncommitted changes
statboard --behind         # only repos behind their remote
statboard --dirty --behind # repos that are both
```

Repos that couldn't be checked always stay visible, even when a filter is active.

### Status indicators

| Status    | Meaning                                         |
|-----------|-------------------------------------------------|
| `OK`      | Clean and in sync with its remote               |
| `DIRTY`   | Has staged, modified, or untracked files        |
| `BEHIND`  | Remote has commits you don't have locally       |
| `AHEAD`   | You have local commits that aren't pushed       |
| `UNKNOWN` | No upstream branch, or remote state unavailable |
| `ERROR`   | The repo couldn't be checked (details below the table) |

If a repo has more than one issue, the most urgent status is shown, and the CHANGES and REMOTE columns carry the details.

> **Note:** Ahead/behind counts reflect each repo's last fetch. StatBoard never runs `git fetch`, so run it yourself if you want fresh remote data.

### Configuration

Tracked repos are stored in a YAML file created by `statboard config add`:

| OS      | Location                                      |
|---------|-----------------------------------------------|
| Windows | `%AppData%\statboard\config.yaml`             |
| macOS   | `~/Library/Application Support/statboard/config.yaml` |
| Linux   | `~/.config/statboard/config.yaml`             |

Set the `NO_COLOR` environment variable to disable colored output.

## How it works

For each tracked repo, StatBoard runs one command:

```
git --no-optional-locks -C <path> status --porcelain=v2 --branch
```

and parses the machine-readable output for the branch, upstream, ahead/behind counts, and changed files. Checks run in goroutines with a concurrency limit, and each writes to its own slot in a results slice, so the table order stays stable.

## Project structure

```
statboard/
├── main.go
├── cmd/                  # Cobra commands (root dashboard, config add/remove/list)
└── internal/
    ├── config/           # YAML config load/save
    ├── git/              # git status parsing and concurrent checks
    ├── filter/           # --dirty / --behind filtering
    └── render/           # table and summary output
```

## Development

```bash
go test ./...   # run unit tests
go vet ./...    # static checks
```

## Built with

[Go](https://go.dev/) · [Cobra](https://github.com/spf13/cobra) · [Lip Gloss](https://github.com/charmbracelet/lipgloss) · [yaml.v3](https://github.com/go-yaml/yaml)

## Author

Allen Baez, Computer Science student at Montclair State University
