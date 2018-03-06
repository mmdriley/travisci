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

	repos, err := client.RepositoriesForCurrentUser()
	if err != nil {
		fmt.Printf("ERROR: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("%+v\n", repos.Pagination.First)
	fmt.Printf("%+v\n", repos.Pagination.Next)
	fmt.Printf("%+v\n", repos.Pagination.Last)

	for _, repo := range repos.Repositories {
		fmt.Printf("%s\n", repo.Slug)
	}

	// fmt.Printf("%+v\n", repos.Repositories)
	// fmt.Printf("%+v\n", repos.Pagination)

	envvars, err := client.EnvVarsForRepository("pulumi/pulumi")
	if err != nil {
		fmt.Printf("ERROR: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("%+v\n", envvars)
}
