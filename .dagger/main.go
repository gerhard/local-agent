package main

import (
	"context"

	"dagger/local-agent/internal/dagger"

	"github.com/google/uuid"
)

type LocalAgent struct{}

// Creates a development environment for the project, installs all the needed tools and libraries
func (l *LocalAgent) DevEnvironment(
	ctx context.Context,
	// Codebase to work on
	source *dagger.Directory,
) (*dagger.Container, error) {
	// Reads the system prompt file from the current module
	systemPrompt, err := dag.CurrentModule().Source().File("qwen_system_prompt.md").Contents(ctx)
	if err != nil {
		return nil, err
	}

	// Create an environment around the source directory:
	// - an alpine based container with the source directory mounted
	// - a set of tools available to the LLM to read files, install packages, etc.
	env := dag.AlpineWorkspace(source)

	return dag.LLM().
		WithAlpineWorkspace(env).

		WithSystemPrompt(systemPrompt).
		WithPromptVar("assignment", "Create a development environment for the project, install all the needed tools and libraries").
		WithPromptFile(dag.CurrentModule().Source().File("qwen_dev_env.md")).

		AlpineWorkspace().Container(), nil
}

// Enter a development environment and export the workspace directory
func (l *LocalAgent) WorkOn(
	ctx context.Context,
	// Codebase to work on
	source *dagger.Directory,
) (*dagger.Directory, error) {
	devEnv, err := l.DevEnvironment(ctx, source)
	if err != nil {
		return nil, err
	}
	return devEnv.
		WithMountedCache("/terminal", dag.CacheVolume(uuid.New().String())).
		WithExec([]string{"cp", "-r", "/workspace", "/terminal"}).
		WithWorkdir("/terminal/workspace").
		Terminal().
		WithExec([]string{"cp", "-r", "/terminal/workspace", "/out"}).
		Directory("/out"), nil
}
