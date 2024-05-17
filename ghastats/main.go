// gh module + https://github.com/fchimpan/gh-workflow-stats

package main

import "context"

type Ghastats struct{}

// Stats for a specific GitHub Actions Workflow
func (m *Ghastats) Run(
	ctx context.Context,

	// GitHub Token - required by gh cli
	token *Secret,

	// GitHub Org / Owner
	//
	// +optional
	// +default="dagger"
	org string,

	// GitHub Repository
	//
	// +optional
	// +default="dagger"
	repo string,

	// GitHub Actions Workflow path
	//
	// +optional
	// +default="engine-and-cli.yml"
	workflow string,

	// +optional
	// +default=">2024-05-15"
	// TODO: default to last 7 days
	since string,

) (string, error) {
	return dag.Gh().WithToken(token).
		Run("extensions install fchimpan/gh-workflow-stats").
		WithExec([]string{"gh", "workflow-stats", "--org", org, "--repo", repo, "--file", workflow, "--branch", "main", "--created", since}).
		Stdout(ctx)
}
