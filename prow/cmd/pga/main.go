package main

import (
	"fmt"
	"http"
	"net/http"
	"os"
)

func main() {
	// Consume GHA runtime info via env vars
	// ref https://docs.github.com/en/github-ae@latest/actions/learn-github-actions/environment-variables#default-environment-variables
	// GITHUB_EVENT_NAME, GITHUB_EVENT_PATH
	ghEvent := os.Getenv(GITHUB_EVENT_NAME)
	ghEventPayload := os.Getenv(GITHUB_EVENT_PATH)

	// Let's go directly to the fn handler as per @hh suggestion
	//
}
