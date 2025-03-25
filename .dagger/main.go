package main

import (
	"dagger/local-agent/internal/dagger"

	"github.com/google/uuid"
)

type LocalAgent struct{}

// Creates a development environment for the project, installs all the needed tools and libraries
func (l *LocalAgent) DevEnvironment(
	// Codebase to work on
	source *dagger.Directory,
) *dagger.Container {
	return dag.
		LLM().
		WithPromptVar("assignment", "Create a development environment for the project, install all the needed tools and libraries").
		WithPromptFile(dag.CurrentModule().Source().File("qwen_dev_env.md")).
		WithDevWorkspace(dag.DevWorkspace(source)).
		DevWorkspace().
		Container()
}

// Enter a development environment and export the workspace directory
func (l *LocalAgent) WorkOn(
	// Codebase to work on
	source *dagger.Directory,
) *dagger.Directory {
	return l.DevEnvironment(source).
		WithMountedCache("/terminal", dag.CacheVolume(uuid.New().String())).
		WithExec([]string{"cp", "-r", "/workspace", "/terminal"}).
		WithWorkdir("/terminal/workspace").
		Terminal().
		WithExec([]string{"cp", "-r", "/terminal/workspace", "/out"}).
		Directory("/out")
}
