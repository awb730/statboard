package git

import "testing"

func TestParse(t *testing.T) {
	tests := []struct {
		name   string
		output string
		want   RepoStatus
	}{
		{
			name:   "clean and in sync",
			output: "# branch.oid abc123\n# branch.head main\n# branch.upstream origin/main\n# branch.ab +0 -0\n",
			want:   RepoStatus{Branch: "main", Upstream: "origin/main", RemoteKnown: true},
		},
		{
			name:   "ahead and behind",
			output: "# branch.head main\n# branch.upstream origin/main\n# branch.ab +2 -1\n",
			want:   RepoStatus{Branch: "main", Upstream: "origin/main", RemoteKnown: true, Ahead: 2, Behind: 1},
		},
		{
			name: "staged, modified and untracked",
			output: "# branch.head main\n# branch.upstream origin/main\n# branch.ab +0 -0\n" +
				"1 M. N... 100644 100644 100644 aaa bbb staged.go\n" +
				"1 .M N... 100644 100644 100644 aaa bbb modified.go\n" +
				"? new.txt\n",
			want: RepoStatus{Branch: "main", Upstream: "origin/main", RemoteKnown: true, Staged: 1, Modified: 1, Untracked: 1},
		},
		{
			name:   "no upstream means remote unknown",
			output: "# branch.oid abc123\n# branch.head feature\n",
			want:   RepoStatus{Branch: "feature"},
		},
		{
			name:   "windows line endings",
			output: "# branch.head main\r\n# branch.upstream origin/main\r\n# branch.ab +1 -0\r\n? new.txt\r\n",
			want:   RepoStatus{Branch: "main", Upstream: "origin/main", RemoteKnown: true, Ahead: 1, Untracked: 1},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got RepoStatus
			parse(&got, tt.output)
			if got != tt.want {
				t.Errorf("got %+v, want %+v", got, tt.want)
			}
		})
	}
}
