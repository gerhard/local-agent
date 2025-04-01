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
	// Create an environment around the source directory:
	env := dag.Env().
		// an alpine based workspace in input, containing:
		// - an alpine based container with the source directory mounted as the input
		// - a set of tools available to the LLM to read files, install packages, etc.
		WithAlpineWorkspaceInput("workspace", dag.AlpineWorkspace(source), "Alpine Environment with the source codebase mounted").
		// a development environment in output, with all the tools installed
		WithAlpineWorkspaceOutput("result", "Workspace with the development tools installed")

	return dag.LLM().
		WithEnv(env).
		WithPromptFile(dag.CurrentModule().Source().File("qwen_dev_env.md")).
		Env().Output("result").AsAlpineWorkspace().Container()
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
