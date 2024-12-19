package main

import (
	"context"
	"dagger/dagrr/internal/dagger"
	"fmt"
)

type DagrFly struct {
	// +private
	Dagr *Dagr
	// +private
	Flyio *dagger.Flyio
	// +private
	PrimaryRegion string
	// +private
	Size string
	// +private
	Memory string
	// +private
	Disk string
	// +private
	GPU string
	// +private
	Env []string
}

// App manifest: `dagger call on-flyio --token=env:FLY_API_TOKEN manifest file --path=fly.toml export --path=fly.toml`
func (m *DagrFly) Manifest() *dagger.Directory {
	var primaryRegion string
	if m.PrimaryRegion != "" {
		// workaround to leave the config untouched if the region isn't set
		primaryRegion = fmt.Sprintf("primary_region = %q", m.PrimaryRegion)
	}

	engineImageFlavor := ""
	if m.GPU != "" {
		engineImageFlavor = "-gpu"
	}

	envConfig := ""
	for _, envVar := range m.Env {
		if envConfig == "" {
			envConfig = "[env]\n"
		}
		envConfig = fmt.Sprintf("%s  %s\n", envConfig, envVar)
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
`, m.Dagr.App, primaryRegion, envConfig, m.Dagr.Version, engineImageFlavor, m.Disk, m.Size)

	if m.Memory != "" {
		toml = fmt.Sprintf("%s  memory = %q\n", toml, m.Memory)
	}

	if m.GPU != "" {
		toml = fmt.Sprintf("%s  gpu_kind = %q\n", toml, m.GPU)
	}

	return dag.Directory().WithNewFile("fly.toml", toml)
}

// Deploy with default manifest: `dagger call on-flyio --token=env:FLY_API_TOKEN deploy`
// Then: `export _EXPERIMENTAL_DAGGER_RUNNER_HOST=tcp://<APP_NAME>.internal:2345`
// Assumes https://fly.io/docs/networking/private-networking (clashes with Tailscale MagicDNS)
func (m *DagrFly) Deploy(
	ctx context.Context,
	// +optional
	dir *dagger.Directory,
	// +optional
	regions []string,
) (string, error) {
	create, err := m.Flyio.Create(ctx, m.Dagr.App)
	if err != nil {
		return create, err
	}

	if dir == nil {
		dir = m.Manifest()
	}

	return m.Flyio.Deploy(ctx, dir, dagger.FlyioDeployOpts{
		Regions: regions,
	})
}

// Destroy the application: `dagger call on-flyio --token=env:FLY_API_TOKEN destroy`
func (m *DagrFly) Destroy(ctx context.Context) {
	m.Flyio.Destroy(ctx, m.Dagr.App)
}
