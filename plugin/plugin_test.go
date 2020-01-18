// Copyright 2019 the Drone Authors. All rights reserved.
// Use of this source code is governed by the Polyform License
// that can be found in the LICENSE file.

package plugin

import "testing"

func TestPlugin(t *testing.T) {
	t.Skip()
}

func TestMatch(t *testing.T) {
	tests := []struct {
		match    string
		patterns []string
		matched  bool
	}{
		// event matching
		{"push", []string{"pull_request", "push"}, true},
		{"tag", []string{"pull_request", "push"}, false},
		{"tag", []string{}, true},

		// repo matching
		{"octocat/hello-world", []string{"octocat/hello-world", "octocat/spoon-fork"}, true},
		{"octocat/hello-world", []string{"octocat/*", "octocat/spoon-fork"}, true},
		{"octocat/hello-world", []string{}, true},
		{"octocat/hello-world", []string{"spaceghost/*"}, false},
		{"octocat/hello-world", []string{"spaceghost/hello-world"}, false},
	}

	for _, test := range tests {
		got, want := match(test.match, test.patterns), test.matched
		if got != want {
			t.Errorf("Want matched %v for string %q and patterns %v",
				test.matched,
				test.match,
				test.patterns,
			)
		}
	}
}
