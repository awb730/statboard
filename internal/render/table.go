package render

import (
	"fmt"
	"io"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"

	"statboard/internal/git"
)

const statusCol = 5 // index of the STATUS column

// dottedBorder uses dashed box-drawing characters for the lines.
var dottedBorder = lipgloss.Border{
	Top:          "┄",
	Bottom:       "┄",
	Left:         "┆",
	Right:        "┆",
	TopLeft:      "┌",
	TopRight:     "┐",
	BottomLeft:   "└",
	BottomRight:  "┘",
	MiddleLeft:   "├",
	MiddleRight:  "┤",
	Middle:       "┼",
	MiddleTop:    "┬",
	MiddleBottom: "┴",
}

// statusCell picks the most urgent indicator and its color.
func statusCell(s git.RepoStatus) (string, lipgloss.Color) {
	switch {
	case s.Err != nil:
		return "ERROR", lipgloss.Color("9")
	case s.IsDirty():
		return "DIRTY", lipgloss.Color("11")
	case s.IsBehind():
		return "BEHIND", lipgloss.Color("11")
	case s.RemoteKnown && s.Ahead > 0:
		return "AHEAD", lipgloss.Color("14")
	case !s.RemoteKnown:
		return "UNKNOWN", lipgloss.Color("13")
	default:
		return "OK", lipgloss.Color("10")
	}
}

// Table prints the repos as a bordered grid.
func Table(w io.Writer, statuses []git.RepoStatus) {
	rows := make([][]string, 0, len(statuses))
	for _, s := range statuses {
		label, _ := statusCell(s)
		if s.Err != nil {
			rows = append(rows, []string{s.Name, s.Path, "-", "-", "-", label})
			continue
		}
		rows = append(rows, []string{s.Name, s.Path, s.Branch, changes(s), remote(s), label})
	}

	t := table.New().
		Border(dottedBorder).
		BorderStyle(lipgloss.NewStyle().Foreground(lipgloss.Color("8"))).
		BorderRow(true). // dotted line between every row; remove for a compact table
		Headers("REPO", "PATH", "BRANCH", "CHANGES", "REMOTE", "STATUS").
		Rows(rows...).
		StyleFunc(func(row, col int) lipgloss.Style {
			base := lipgloss.NewStyle().Padding(0, 1)
			switch {
			case row == table.HeaderRow:
				return base.Bold(true)
			case col == statusCol:
				_, c := statusCell(statuses[row])
				return base.Foreground(c).Bold(true)
			}
			return base
		})

	fmt.Fprintln(w, t.Render())
}

// Summary prints any error details, the counts, and the remote-freshness note.
func Summary(w io.Writer, shown []git.RepoStatus, total int) {
	errCount := 0
	for _, s := range shown {
		if s.Err != nil {
			errCount++
		}
	}
	if errCount > 0 {
		fmt.Fprintln(w, "\nErrors:")
		for _, s := range shown {
			if s.Err != nil {
				fmt.Fprintf(w, "  %s: %s\n", s.Name, oneLine(s.Err.Error()))
			}
		}
	}

	fmt.Fprintf(w, "\nChecked %d repo(s)", total)
	if len(shown) != total {
		fmt.Fprintf(w, ", showing %d (filtered)", len(shown))
	}
	if errCount > 0 {
		fmt.Fprintf(w, ", %d with errors", errCount)
	}
	fmt.Fprintln(w)
	fmt.Fprintln(w, "Remote counts reflect each repo's last fetch. StatBoard never fetches or changes repos.")
}

func changes(s git.RepoStatus) string {
	var parts []string
	if s.Staged > 0 {
		parts = append(parts, fmt.Sprintf("%d staged", s.Staged))
	}
	if s.Modified > 0 {
		parts = append(parts, fmt.Sprintf("%d modified", s.Modified))
	}
	if s.Untracked > 0 {
		parts = append(parts, fmt.Sprintf("%d untracked", s.Untracked))
	}
	if len(parts) == 0 {
		return "clean"
	}
	return strings.Join(parts, ", ")
}

func remote(s git.RepoStatus) string {
	switch {
	case !s.RemoteKnown && s.Upstream == "":
		return "no upstream"
	case !s.RemoteKnown:
		return "unknown (upstream gone)"
	case s.Ahead == 0 && s.Behind == 0:
		return "in sync"
	case s.Behind == 0:
		return fmt.Sprintf("%d ahead", s.Ahead)
	case s.Ahead == 0:
		return fmt.Sprintf("%d behind", s.Behind)
	default:
		return fmt.Sprintf("%d ahead, %d behind", s.Ahead, s.Behind)
	}
}

// oneLine keeps multi-line git errors from breaking the table layout.
func oneLine(msg string) string {
	return strings.ReplaceAll(strings.TrimSpace(msg), "\n", " ")
}



