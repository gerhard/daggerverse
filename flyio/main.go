// Manage apps on https://fly.io
//
// Currently only deploys.
// Assumes there is a valid fly.toml in the path provided to --dir.
package main

import (
	"context"
	"fmt"
)

const (
	// https://hub.docker.com/r/flyio/flyctl/tags
	version = "0.2.13"
)

type Flyio struct{}

// Example usage: "dagger call deploy --dir . --token env:FLY_API_TOKEN"
func (m *Flyio) Deploy(ctx context.Context, dir *Directory, token *Secret) (string, error) {
	return dag.Container().
		From(fmt.Sprintf("flyio/flyctl:v%s", version)).
		WithMountedDirectory("/app", dir).
		WithWorkdir("/app").
		WithSecretVariable("FLY_API_TOKEN", token).
		WithExec([]string{"deploy"}).
		Stdout(ctx)
}
