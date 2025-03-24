package main

import (
	"context"

	"dagger/local-agent/internal/dagger"
)

type LocalAgent struct{}

// Summarize a subreddit
func (l *LocalAgent) Summarize(
	ctx context.Context,
	clientId, clientSecret, username, password *dagger.Secret,
	// Name of the subreddit to summarize
	subreddit string,
	// Include comments or only summarize posts
	// +optional
	comments bool,
) (string, error) {
	return dag.
		LLM(dagger.LLMOpts{Model: "ignaciolopezluna020/llama3.2:1B"}).
		SetReddit("reddit_fetcher", dag.Reddit(clientId, clientSecret, username, password)).
		WithPrompt("Create a summary of the subreddit '" + subreddit + "'").
		WithPrompt("You have access to a reddit fetcher to get recent posts of the subreddit").
		With(func(r *dagger.LLM) *dagger.LLM {
			if comments {
				return r.WithPrompt("Use the reddit fetcher to get recent comments of the subreddit")
			}
			return r
		}).
		WithPrompt("Write a few sentences about the subreddit '" + subreddit + "' and highlight the most interesting posts.").
		WithPrompt("Include some exceprt from the posts").
		WithPrompt("Include some statistics about the subreddit").
		WithPrompt("Format the response in markdown").
		LastReply(ctx)
}

// Print a summary of the subreddit
func (l *LocalAgent) PrintSummary(
	ctx context.Context,
	clientId, clientSecret, username, password *dagger.Secret,
	// Name of the subreddit to summarize
	subreddit string,
	// Include comments or only summarize posts
	// +optional
	comments bool,
) (string, error) {
	summary, err := l.Summarize(ctx, clientId, clientSecret, username, password, subreddit, comments)
	if err != nil {
		return "", err
	}
	return dag.Glow().DisplayMarkdown(ctx, summary)
}

// Creates a development environment for the project, installs all the needed tools and libraries
func (l *LocalAgent) DevEnvironment(
	source *dagger.Directory,
) *dagger.Container {
	return dag.
		LLM(dagger.LLMOpts{Model: "eunomie/qwen2.5-coder-14b-instruct:q5_k_m"}).
		WithPromptVar("assignment", "Create a development environment for the project, install all the needed tools and libraries").
		WithPromptFile(dag.CurrentModule().Source().File("qwen_dev_env.md")).
		WithDevWorkspace(dag.DevWorkspace(source)).
		DevWorkspace().
		Container()
}
