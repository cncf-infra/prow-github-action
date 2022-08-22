package main

import (
	"fmt"
	"os"
)

func main() {
	// Consume GHA runtime info via env vars
	// ref https://docs.github.com/en/github-ae@latest/actions/learn-github-actions/environment-variables#default-environment-variables
	// GITHUB_EVENT_NAME, GITHUB_EVENT_PATH
	ghEvent := os.Getenv("GITHUB_EVENT_NAME")
	ghEventPayload := os.Getenv("GITHUB_EVENT_PATH")

	fmt.Printf("event,payload %s: %s\n", ghEvent, ghEventPayload)
	//Load a handler
	Handler.routeEventToHndler(ghEvent, "eventGUID", []byte(ghEventPayload))
}
