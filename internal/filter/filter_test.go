package filter

import (
	"errors"
	"testing"

	"statboard/internal/git"
)

func TestApplyKeepsErrorRows(t *testing.T) {
	statuses := []git.RepoStatus{
		{Name: "clean"},
		{Name: "dirty", Modified: 1},
		{Name: "broken", Err: errors.New("boom")},
	}

	got := Apply(statuses, Options{Dirty: true})

	if len(got) != 2 || got[0].Name != "dirty" || got[1].Name != "broken" {
		t.Errorf("expected [dirty broken], got %+v", got)
	}
}