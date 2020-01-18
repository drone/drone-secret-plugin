// Copyright 2019 the Drone Authors. All rights reserved.
// Use of this source code is governed by the Polyform License
// that can be found in the LICENSE file.

package plugin

import (
	"context"
	"testing"

	"github.com/drone/drone-go/drone"
	"github.com/drone/drone-go/plugin/secret"

	"github.com/google/go-cmp/cmp"
)

var noContext = context.Background()

func TestPlugin(t *testing.T) {
	secrets := []*Secret{
		{
			Name:   "username",
			Value:  "root",
			Repos:  []string{},
			Events: []string{},
		},
		{
			Name:   "password",
			Value:  "correct-horse-battery-staple",
			Repos:  []string{},
			Events: []string{},
		},
	}
	req := &secret.Request{
		Name:  "username",
		Repo:  drone.Repo{},
		Build: drone.Build{},
	}
	plugin := New(secrets)
	got, err := plugin.Find(noContext, req)
	if err != nil {
		t.Error(err)
		return
	}
	want := &drone.Secret{
		Name: "username",
		Data: "root",
		Pull: true,
		Fork: true,
	}
	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("Unexpected secret")
		t.Logf(diff)
	}
}

func TestPlugin_Match(t *testing.T) {
	secrets := []*Secret{
		{
			Name:   "username",
			Value:  "root",
			Repos:  []string{"octocat/*"},
			Events: []string{},
		},
	}
	req := &secret.Request{
		Name:  "username",
		Repo:  drone.Repo{Slug: "octocat/hello-world"},
		Build: drone.Build{Event: "push"},
	}
	plugin := New(secrets)
	got, err := plugin.Find(noContext, req)
	if err != nil {
		t.Error(err)
		return
	}
	want := &drone.Secret{
		Name: "username",
		Data: "root",
		Pull: true,
		Fork: true,
	}
	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("Unexpected secret")
		t.Logf(diff)
	}
}

func TestPlugin_NoMatch(t *testing.T) {
	secrets := []*Secret{
		{
			Name:  "username",
			Value: "root",
		},
	}
	req := &secret.Request{
		Name:  "password",
		Repo:  drone.Repo{Slug: "octocat/hello-world"},
		Build: drone.Build{Event: "push"},
	}
	plugin := New(secrets)
	got, err := plugin.Find(noContext, req)
	if err != nil {
		t.Error(err)
		return
	}
	if got != nil {
		t.Errorf("Expect nil secret due to event mismatch")
	}
}

func TestPlugin_NoMatch_Repo(t *testing.T) {
	secrets := []*Secret{
		{
			Name:   "username",
			Value:  "root",
			Repos:  []string{"octocat/*"},
			Events: []string{},
		},
	}
	req := &secret.Request{
		Name:  "username",
		Repo:  drone.Repo{Slug: "spaceghost/hello-world"},
		Build: drone.Build{Event: "push"},
	}
	plugin := New(secrets)
	got, err := plugin.Find(noContext, req)
	if err != nil {
		t.Error(err)
		return
	}
	if got != nil {
		t.Errorf("Expect nil secret due to repository mismatch")
	}
}

func TestPlugin_NoMatch_Event(t *testing.T) {
	secrets := []*Secret{
		{
			Name:   "username",
			Value:  "root",
			Repos:  []string{"octocat/*"},
			Events: []string{"pull_request"},
		},
	}
	req := &secret.Request{
		Name:  "username",
		Repo:  drone.Repo{Slug: "octocat/hello-world"},
		Build: drone.Build{Event: "push"},
	}
	plugin := New(secrets)
	got, err := plugin.Find(noContext, req)
	if err != nil {
		t.Error(err)
		return
	}
	if got != nil {
		t.Errorf("Expect nil secret due to event mismatch")
	}
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
