package travisci

import (
	"os"

	"github.com/pkg/errors"
)

// https://dave.cheney.net/2014/10/17/functional-options-for-friendly-apis

/*

recall our options progression is:

endpoint:
1. passed to constructor
2. TRAVIS_ENDPOINT
3. default_endpoint in config
4. travis.org

token:
1. passed to constructor as Travis or GitHub token
2. TRAVIS_TOKEN
3. access_token in config -- *depends on endpoint!*

TODO for client:
-H "Travis-API-Version: 3"
-H "Content-Type: application/json"

*/

var (
	orgURI = "https://api.travis-ci.org/"
	proURI = "https://api.travis-ci.com/"
)

type Client struct {
	endpoint string

	// If githubToken is set and token isn't, we'll hit endpoint to exchange.
	githubToken string
	token       string
}

type Option func(*Client)

func AccessToken(token string) Option {
	return func(c *Client) {
		c.token = token
	}
}

func GitHubAccessToken(token string) Option {
	return func(c *Client) {
		c.githubToken = token

		// Set token so we know it's configured. This value will be overwritten.
		c.token = "github"
	}
}

func Endpoint(endpoint string) Option {
	return func(c *Client) {
		c.endpoint = endpoint
	}
}

func OrgEndpoint() Option {
	return Endpoint(orgURI)
}

func ProEndpoint() Option {
	return Endpoint(proURI)
}

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

// Set token and endpoint the way travis.rb would:
// Endpoint: first env var, then config, then fall back to travis-ci.org
// Token: first env var, then config for chosen endpoint
func (c *Client) configureTokenAndEndpoint() error {
	if c.token != "" && c.endpoint != "" {
		// Already configured explicitly.
		return nil
	}

	// Try environment variables.
	if c.token == "" {
		c.token = os.Getenv("TRAVIS_TOKEN")
	}
	if c.endpoint == "" {
		c.endpoint = os.Getenv("TRAVIS_ENDPOINT")
	}

	// Exit now if we're done and don't need to read config.
	if c.token != "" && c.endpoint != "" {
		return nil
	}

	// Read config to look for default endpoint and saved token.
	config, err := readConfig()
	if err != nil {
		return errors.Wrap(err, "reading Travis config file")
	}

	if c.endpoint == "" {
		c.endpoint = config.DefaultEndpoint
	}

	// If we don't have an endpoint by now, fall back to travis-ci.org.
	// We need to choose an endpoint before looking for a saved token.
	if c.endpoint == "" {
		c.endpoint = orgURI
	}

	if c.token == "" {
		if endpointConfig, ok := config.Endpoints[c.endpoint]; ok {
			c.endpoint = endpointConfig.AccessToken
		}
	}

	if c.token == "" {
		return errors.New("no access token provided")
	}

	return nil
}

func NewClient(options ...Option) interface{} {
	c := Client{}
	for _, option := range options {
		option(&c)
	}

	var err error
	if err = c.configureTokenAndEndpoint(); err != nil {
		return err
	}

	// If we have a GitHub access token, get a Travis token from it.
	if c.githubToken != "" {
		c.token, err = travisTokenFromGitHubToken(c.githubToken, c.endpoint)
		if err != nil {
			return err
		}

		c.githubToken = "" // No need to keep this around.
	}

	return c
}
