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
		WithAlpineWorkspaceInput(
			"workspace",
			dag.AlpineWorkspace(source),
			"An alpine workspace containing the source code directory.").
		// a development environment in output, with all the tools installed
		WithAlpineWorkspaceOutput(
			"result",
			"The updated alpine workspace with the necessary development tools and project dependencies installed, based on the analyzed source directory")

	return dag.LLM().
		WithEnv(env).
		WithPrompt("do what you need to do").
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
