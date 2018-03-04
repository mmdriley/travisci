package main

import (
	"fmt"
	"os"

	"github.com/mmdriley/travisci"
)

func main() {
	// client, err := travisci.NewClient(travisci.GitHubAccessToken(os.Getenv("GITHUB_TOKEN")))
	client, err := travisci.NewClient()
	if err != nil {
		fmt.Printf("ERROR: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("%+v\n", client)

	user, err := client.CurrentUser()
	if err != nil {
		fmt.Printf("ERROR: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("%+v\n", user)
}
