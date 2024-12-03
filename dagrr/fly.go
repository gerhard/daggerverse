package main

import (
	"context"
	"dagger/dagrr/internal/dagger"
	"fmt"
)

type DagrrFly struct {
	// +private
	Dagrr *Dagrr
	// +private
	Flyio *dagger.Flyio
}

// App manifest: `dagger call on-flyio --token=env:FLY_API_TOKEN manifest file --path=fly.toml export --path=fly.toml`
func (m *DagrrFly) Manifest(
	// Primary region to use for deploying new machines, see https://fly.io/docs/reference/configuration/#primary-region
	// +optional
	primaryRegion string,

	// VM size, see https://fly.io/docs/about/pricing/#compute
	//
	// +optional
	// +default="performance-2x"
	size string,

	// Memory to request, see https://fly.io/docs/reference/configuration/#memory
	// +optional
	memory string,

	// GPU kind to use, see https://fly.io/docs/reference/configuration/#gpu_kind
	// +optional
	gpuKind string,

	// Disk size in GB
	//
	// +optional
	// +default="100GB"
	disk string,

	// Environment variables to export on the machine, see https://fly.io/docs/reference/configuration/#the-env-variables-section
	// Each env var needs to follow the TOML format (eg. MY_KEY = "value")
	// FIXME(samalba): turn this into a map[string]string once supported by the Dagger Go SDK
	// +optional
	environment []string,
) *dagger.Directory {
	if primaryRegion != "" {
		// workaround to leave the config untouched if the region isn't set
		primaryRegion = fmt.Sprintf("primary_region = %q", primaryRegion)
	}

	envVars := ""
	for _, envVar := range environment {
		if envVars == "" {
			envVars = "[env]\n"
		}
		envVars = fmt.Sprintf("%s  %s\n", envVars, envVar)
	}

	engineImageFlavor := ""
	if gpuKind != "" {
		engineImageFlavor = "-gpu"
	}

	toml := fmt.Sprintf(`# https://fly.io/docs/reference/configuration/

app = "%s"
%s

kill_signal = "SIGINT"
kill_timeout = 30

%s

[build]
  image = "registry.dagger.io/engine:v%s%s"

[mounts]
  source = "dagger"
  destination = "/var/lib/dagger"
  initial_size = "%s"

[processes]
  dagger = "--addr unix:///var/run/buildkit/buildkitd.sock --addr tcp://0.0.0.0:2345"

[checks]
  [checks.http]
    grace_period = "3s"
    interval = "2s"
    port = 2345
    timeout = "1s"
    type = "tcp"

[[services]]
  internal_port = 2345
  protocol = "tcp"
  auto_stop_machines = false
  auto_start_machines = true
  min_machines_running = 1
  processes = ["dagger"]

  [[services.ports]]
    handlers = ["http"]
    port = 2345

[[vm]]
  size = "%s"
`, m.Dagrr.App, primaryRegion, envVars, m.Dagrr.Version, engineImageFlavor, disk, size)

	if memory != "" {
		toml = fmt.Sprintf("%s  memory = %q\n", toml, memory)
	}

	if gpuKind != "" {
		toml = fmt.Sprintf("%s  gpu_kind = %q\n", toml, gpuKind)
	}

	return dag.Directory().WithNewFile("fly.toml", toml)
}

// Deploy with default manifest: `dagger call on-flyio --token=env:FLY_API_TOKEN deploy`
// Then: `export _EXPERIMENTAL_DAGGER_RUNNER_HOST=tcp://<APP_NAME>.internal:2345`
// Assumes https://fly.io/docs/networking/private-networking (clashes with Tailscale MagicDNS)
func (m *DagrrFly) Deploy(
	ctx context.Context,
	// +optional
	dir *dagger.Directory,
	// +optional
	regions []string,
) (string, error) {
	create, err := m.Flyio.Create(ctx, m.Dagrr.App)
	if err != nil {
		return create, err
	}

	if dir == nil {
		dir = m.Manifest("100GB", "performance-2x", "", "", "", nil)
	}

	return m.Flyio.Deploy(ctx, dir, dagger.FlyioDeployOpts{
		Regions: regions,
	})
}

// Destroy the application: `dagger call on-flyio --token=env:FLY_API_TOKEN destroy`
func (m *DagrrFly) Destroy(ctx context.Context) {
	m.Flyio.Destroy(ctx, m.Dagrr.App)
}
