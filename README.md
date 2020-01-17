This plugin provides global secrets for use in your pipelines. This plugin is a direct port of the global secret file in Drone 0.8. _Please note this project requires Drone server version 1.4 or higher._

## Secret File

Secrets are loaded from a yaml configuration file. Example secrets configuration file:

```text
- name: docker_username
  value: octocat

- name: docker_password
  value: correct-horse-battery-staple
  repos: [ octocat/hello-world, github/* ]
  events: [ push, tag ]
```

## Installation

Create a shared secret:

```console
$ openssl rand -hex 16
bea26a2221fd8090ea38720fc445eca6
```

Download and run the plugin:

```console
$ docker run -d \
  --publish=3000:3000 \
  --env=DRONE_DEBUG=true \
  --env=DRONE_SECRET=bea26a2221fd8090ea38720fc445eca6 \
  --env=DRONE_SECRET_FILE=/etc/secrets.yml \
  --restart=always \
  --volume=/etc/secrets.yml:/etc/secrets.yml \
  --name=secrets drone/secret-plugin
```

Update your Drone server configuration to include the plugin address and the shared secret.

```text
DRONE_SECRET_PLUGIN_ENDPOINT=http://1.2.3.4:3000
DRONE_SECRET_PLUGIN_SECRET=bea26a2221fd8090ea38720fc445eca6
