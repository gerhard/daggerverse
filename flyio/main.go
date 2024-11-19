// Manage apps on https://fly.io
//
// Currently only deploys.
// Assumes there is a valid fly.toml in the path provided to --dir.
package main

import (
	"context"
	"fmt"
	"main/internal/dagger"
	"time"
)

const (
	// https://hub.docker.com/r/flyio/flyctl/tags
	latestVersion = "0.2.93"
)

type Flyio struct {
	// +private
	Container *dagger.Container
	// +private
	Version string
	// +private
	Org string
}

func New(
	// fly auth token: `--token=env:FLY_API_TOKEN`
	token *dagger.Secret,

	// Fly.io org where all operations will run in, defaults to: `--org=personal`
	//
	// +optional
	// +default="personal"
	org string,

	// flyctl version to use: `--version=0.2.79`
	//
	// +optional
	version string,

	// Custom container to use as the base container
	//
	// +optional
	container *dagger.Container,
) *Flyio {
	if container == nil {
		if version == "" {
			version = latestVersion
		}

		container = dag.Container().
			From(fmt.Sprintf("flyio/flyctl:v%s", version)).
			WithSecretVariable("FLY_API_TOKEN", token).
			WithEnvVariable("CACHE_BUSTED_AT", time.Now().String())
	}

	return &Flyio{
		Container: container,
		Version:   version,
		Org:       org,
	}
}

// Deploys app from current directory: `dagger call ... deploy --dir=hello-static`
func (m *Flyio) Deploy(
	ctx context.Context,
	// App directory - must contain `fly.toml`
	dir *dagger.Directory,
) (string, error) {
	return m.Container.
		WithMountedDirectory("/app", dir).
		WithWorkdir("/app").
		WithExec([]string{"/flyctl", "deploy"}).
		Stdout(ctx)
}

// Creates app - required for deploy to work: `dagger call ... create --app=gostatic-example-2024-07-03`
func (m *Flyio) Create(
	ctx context.Context,
	// App name: `--app=myapp-2024-07-03`
	app string,

) (string, error) {
	return m.Container.
		WithExec([]string{"/flyctl", "apps", "create", app, "--org", m.Org}).
		Stdout(ctx)
}

// Opens terminal in this app: `dagger call ... terminal --app=gostatic-example-2024-07-03 --interactive`
func (m *Flyio) Terminal(
	ctx context.Context,
	// App name: `--app=myapp-2024-07-03`
	app string,

) *dagger.Container {
	return m.Container.
		WithExec([]string{"/flyctl", "ssh", "console", "--app", app, "--org", m.Org}).
		Terminal()
}
