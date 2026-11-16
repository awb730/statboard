package git

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type RepoStatus struct {
	Name        string
	Path        string
	Branch      string
	Upstream    string // empty if the branch has no upstream configured
	RemoteKnown bool   // true only if git reported ahead/behind counts
	Ahead       int
	Behind      int
	Modified    int
	Staged      int
	Untracked   int
	Err         error
}

func (r RepoStatus) IsDirty() bool {
	return r.Modified+r.Staged+r.Untracked > 0
}

func (r RepoStatus) IsBehind() bool {
	return r.RemoteKnown && r.Behind > 0
}

// Check runs git status for one repo. Failures are stored in Err, never returned.
func Check(path string) RepoStatus {
	r := RepoStatus{Name: filepath.Base(path), Path: path}

	info, err := os.Stat(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			r.Err = errors.New("path not found (moved or deleted?)")
		} else {
			r.Err = fmt.Errorf("cannot access path: %w", err)
		}
		return r
	}
	if !info.IsDir() {
		r.Err = errors.New("path is not a directory")
		return r
	}

	cmd := exec.Command("git", "--no-optional-locks", "-C", path,
		"status", "--porcelain=v2", "--branch")
	out, err := cmd.Output()
	if err != nil {
		var exitErr *exec.ExitError
		switch {
		case errors.Is(err, exec.ErrNotFound):
			r.Err = errors.New("git is not installed or not on PATH")
		case errors.As(err, &exitErr):
			r.Err = fmt.Errorf("git error: %s", cleanGitMsg(string(exitErr.Stderr)))
		default:
			r.Err = fmt.Errorf("could not run git: %w", err)
		}
		return r
	}

	parse(&r, string(out))
	return r
}

// cleanGitMsg keeps only the first line of git's message, without the "fatal: " prefix.
func cleanGitMsg(s string) string {
	s = strings.TrimPrefix(strings.TrimSpace(s), "fatal: ")
	if i := strings.Index(s, "\n"); i >= 0 {
		s = s[:i]
	}
	return s
}

// parse fills r from porcelain v2 output. Kept separate so it's easy to test.
func parse(r *RepoStatus, output string) {
	for _, line := range strings.Split(output, "\n") {
		line = strings.TrimRight(line, "\r")

		switch {
		case strings.HasPrefix(line, "# branch.head "):
			r.Branch = strings.TrimPrefix(line, "# branch.head ")

		case strings.HasPrefix(line, "# branch.ab "):
			counts := strings.TrimPrefix(line, "# branch.ab ")
			if _, err := fmt.Sscanf(counts, "+%d -%d", &r.Ahead, &r.Behind); err != nil {
				r.Err = fmt.Errorf("could not parse ahead/behind: %w", err)
				return
			}
			r.RemoteKnown = true

		case strings.HasPrefix(line, "1 "), strings.HasPrefix(line, "2 "):
			xy := line[2:4]
			if xy[0] != '.' {
				r.Staged++
			}
			if xy[1] != '.' {
				r.Modified++
			}

		case strings.HasPrefix(line, "u "): // merge conflict
			r.Modified++

		case strings.HasPrefix(line, "? "):
			r.Untracked++

		case strings.HasPrefix(line, "# branch.upstream "):
			r.Upstream = strings.TrimPrefix(line, "# branch.upstream ")
		}
	}
}
