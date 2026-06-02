package dockercomposetooler

import (
	"context"
	"os/exec"

	root "github.com/hanzoai/o11y-foundry/internal/tooler"
)

var _ root.Tooler = (*dockerComposeTooler)(nil)

type dockerComposeTooler struct{}

func New() *dockerComposeTooler {
	return &dockerComposeTooler{}
}

func (tooler *dockerComposeTooler) Name() string {
	return "docker-compose"
}

func (tooler *dockerComposeTooler) Gauge(ctx context.Context) error {
	// Legacy standalone binary.
	if err := root.ExecChecker(ctx, "docker-compose"); err == nil {
		return nil
	}

	if err := root.ExecChecker(ctx, "docker"); err == nil {
		if err := exec.CommandContext(ctx, "docker", "compose", "version").Run(); err == nil {
			return nil
		}
	}

	return errors.Newf(errors.TypeNotFound, "neither 'docker-compose' nor the 'docker compose' plugin is available")
}

func (tooler *dockerComposeTooler) Install(ctx context.Context) error {
	return nil
}
