package filter

import "statboard/internal/git"

type Options struct {
	Dirty  bool
	Behind bool
}

// Apply returns the statuses that match every enabled filter.
// Repos with errors always pass, since their state can't be ruled out.
func Apply(statuses []git.RepoStatus, opts Options) []git.RepoStatus {
	if !opts.Dirty && !opts.Behind {
		return statuses
	}

	out := make([]git.RepoStatus, 0, len(statuses))
	for _, s := range statuses {
		if s.Err != nil {
			out = append(out, s)
			continue
		}
		if opts.Dirty && !s.IsDirty() {
			continue
		}
		if opts.Behind && !s.IsBehind() {
			continue
		}
		out = append(out, s)
	}
	return out
}
