// Copyright 2019 the Drone Authors. All rights reserved.
// Use of this source code is governed by the Polyform License
// that can be found in the LICENSE file.

package main

import (
	"io/ioutil"
	"net/http"

	"github.com/drone/drone-go/plugin/secret"
	"github.com/drone/drone-secret-plugin/plugin"
	"gopkg.in/yaml.v2"

	_ "github.com/joho/godotenv/autoload"
	"github.com/kelseyhightower/envconfig"
	"github.com/sirupsen/logrus"
)

// spec provides the plugin settings.
type spec struct {
	Bind   string `envconfig:"DRONE_BIND"`
	Debug  bool   `envconfig:"DRONE_DEBUG"`
	Secret string `envconfig:"DRONE_SECRET"`
	File   string `envconfig:"DRONE_SECRET_FILE"`
}

func main() {
	spec := new(spec)
	err := envconfig.Process("", spec)
	if err != nil {
		logrus.Fatal(err)
	}

	if spec.Debug {
		logrus.SetLevel(logrus.DebugLevel)
	}
	if spec.Secret == "" {
		logrus.Fatalln("missing secret key")
	}
	if spec.Bind == "" {
		spec.Bind = ":3000"
	}
	if spec.File == "" {
		spec.File = "/etc/secrets.yml"
	}

	raw, err := ioutil.ReadFile(spec.File)
	if err != nil {
		logrus.WithError(err).Fatalln("Cannot read secrets file")
	}

	var secrets []*plugin.Secret
	if err := yaml.Unmarshal(raw, &secrets); err != nil {
		logrus.WithError(err).Fatalln("Cannot unmarshal secrets file")
	}

	handler := secret.Handler(
		spec.Secret,
		plugin.New(secrets),
		logrus.StandardLogger(),
	)

	logrus.Infof("server listening on address %s", spec.Bind)

	http.Handle("/", handler)
	logrus.Fatal(http.ListenAndServe(spec.Bind, nil))
}
