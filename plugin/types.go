// Copyright 2019 the Drone Authors. All rights reserved.
// Use of this source code is governed by the Polyform License
// that can be found in the LICENSE file.

package plugin

// Secret represents a secret and its matching rules that
// determine whether or not a pipeline should be granted
// access to the secret value.
type Secret struct {
	Name   string
	Value  string
	Repos  []string
	Events []string
}
