// Manages Dagger Engines on a bunch of platforms

package main

import (
	"context"
	"dagger/dagrr/internal/dagger"
	"strings"
	"time"

	"github.com/0x6flab/namegenerator"
)

type Dagrr struct {
	// +private
	Version string
	// +private
	App string
}

func New(
	ctx context.Context,

	// Dagger version to use (omit for latest): `--version=0.14.0`
	//
	// +optional
	version string,

	// App name, defaults to version & unique name & date: `--app=dagger-v0-14-0-<GENERATED_NAME>-2024-11-19`
	//
	// +optional
	app string,
) (*Dagrr, error) {
	if version == "" {
		// If version isn't set, assume latest
		v, err := dag.Version(ctx)
		if err != nil {
			return nil, err
		}
		version = v[1:]
	}

	m := &Dagrr{
		Version: version,
	}

	if app == "" {
		app = strings.Join([]string{
			"dagger",
			m.versionUrlized(),
			strings.ToLower(namegenerator.NewGenerator().Generate()),
			time.Now().Format("2006-01-02"),
		}, "-")
	}
	m.App = app

	return m, nil
}

// Manages Dagger on Fly.io: `dager call on-flyio --token=env:FLY_API_TOKEN deploy`
func (m *Dagrr) OnFlyio(
	// `flyctl tokens create deploy` then `--token=env:FLY_API_TOKEN`
	token *dagger.Secret,

	// Fly.io org name
	//
	// +optional
	// +default="personal"
	org string,
) *DagrrFly {
	return &DagrrFly{
		Dagrr: m,
		Flyio: dag.Flyio(token, dagger.FlyioOpts{
			Org: org,
		}),
	}
}

func (m *Dagrr) versionUrlized() string {
	return "v" + strings.ReplaceAll(m.Version, ".", "-")
}

// Returns the app name
func (m *Dagrr) GetApp() string {
	return m.App
}
