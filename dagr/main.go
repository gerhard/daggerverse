// Manages Dagger Engines on a bunch of platforms

package main

import (
	"context"
	"dagger/dagrr/internal/dagger"
	"strings"
	"time"

	"github.com/0x6flab/namegenerator"
)

type Dagr struct {
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
) (*Dagr, error) {
	if version == "" {
		// If version isn't set, assume latest
		v, err := dag.Version(ctx)
		if err != nil {
			return nil, err
		}
		version = v[1:]
	}

	m := &Dagr{
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

// Manages Dagger on Fly.io: `dager call on-flyio --token=env:FLY_API_TOKEN`
func (m *Dagr) OnFlyio(
	// `flyctl tokens create deploy` then `--token=env:FLY_API_TOKEN`
	token *dagger.Secret,

	// Fly.io org name
	//
	// +default="personal"
	org string,

	// Primary region to use for deploying new machines, see https://fly.io/docs/reference/configuration/#primary-region
	// +optional
	primaryRegion string,

	// Persistent disk size in GB
	//
	// +default="100GB"
	disk string,

	// VM size, see https://fly.io/docs/about/pricing/#compute
	//
	// +default="performance-2x"
	size string,

	// Memory to request, see https://fly.io/docs/reference/configuration/#memory
	// +optional
	memory string,

	// GPU kind to use, see https://fly.io/docs/reference/configuration/#gpu_kind
	// +optional
	gpuKind string,

	// Environment variables to export on the machine, see https://fly.io/docs/reference/configuration/#the-env-variables-section
	// Each env var needs to follow the TOML format (eg. MY_KEY = "value")
	// FIXME(samalba): turn this into a map[string]string once supported by the Dagger Go SDK
	// +optional
	env []string,
) *DagrFly {

	return &DagrFly{
		Dagr: m,
		Flyio: dag.Flyio(token, dagger.FlyioOpts{
			Org: org,
		}),
		PrimaryRegion: primaryRegion,
		Disk:          disk,
		Size:          size,
		Memory:        memory,
		GPU:           gpuKind,
		Env:           env,
	}
}

func (m *Dagr) versionUrlized() string {
	return "v" + strings.ReplaceAll(m.Version, ".", "-")
}

// Returns the app name: `dagger call get-app`
func (m *Dagr) GetApp() string {
	return m.App
}
