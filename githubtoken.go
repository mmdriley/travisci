package travisci

type GitHubAuthRequest struct {
	GitHubToken string `json:"github_token"`
}

type GitHubAuthResponse struct {
	AccessToken string `json:"access_token"`
}

// curl https://api.travis-ci.com/v3/ | jq .config.github.scopes
func travisTokenFromGitHubToken(githubAccessToken string, endpoint string) (string, error) {
	// curl -X POST https://api.travis-ci.com/auth/github -d '{"github_token":"'$GH_TOKEN'"}' -H 'Content-Type: application/json'
	request := GitHubAuthRequest{GitHubToken: githubAccessToken}

	return "", nil
}
