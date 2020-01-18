// Copyright 2019 the Drone Authors. All rights reserved.
// Use of this source code is governed by the Polyform License
// that can be found in the LICENSE file.

package plugin

import (
	"context"
	"path/filepath"
	"strings"

	"github.com/drone/drone-go/drone"
	"github.com/drone/drone-go/plugin/secret"
)

// New returns a new secret plugin.
func New(secrets []*Secret) secret.Plugin {
	return &plugin{
		secrets: secrets,
	}
}

type plugin struct {
	secrets []*Secret
}

// Find returns the named plugin from the static list of secrets,
// if an only if the name, repository, branch and event match.
func (p *plugin) Find(ctx context.Context, req *secret.Request) (*drone.Secret, error) {
	for _, secret := range p.secrets {
		// if the secret does not match the name, continue
		if !strings.EqualFold(secret.Name, req.Name) {
			continue
		}

		// if the secret does not match the repository name,
		// build event, or branch, skip.
		if !match(req.Build.Event, secret.Events) {
			continue
		}
		if !match(req.Repo.Slug, secret.Repos) {
			continue
		}

		return &drone.Secret{
			Name: secret.Name,
			Data: secret.Value,
			Pull: true,
			Fork: true,
		}, nil
	}
	return nil, nil
}

// helper function returns true if string s matches a
// string in the index.
func match(s string, patterns []string) bool {
	// if the patterns matching list is empty the
	// string is considered a match.
	if len(patterns) == 0 {
		return true
	}
	// if the string matches any pattern in the
	// list it is considered a match.
	for _, pattern := range patterns {
		matched, _ := filepath.Match(pattern, s)
		if matched {
			return true
		}
	}
	return false
}
