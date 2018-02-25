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

func defaultEndpoint(config config) string {
	// Use $TRAVIS_ENDPOINT if defined
	if endpoint := os.Getenv("TRAVIS_ENDPOINT"); endpoint != "" {
		return endpoint
	}

	// Use default endpoint from config.yml
	if config.DefaultEndpoint != "" {
		return config.DefaultEndpoint
	}

	// Absent any other information, use the .org endpoint.
	return orgURI
}

func storedToken(config config, endpoint string) string {
	// Use $TRAVIS_TOKEN if defined
	if token := os.Getenv("TRAVIS_TOKEN"); token != "" {
		return token
	}

	// Read token for the given endpoint from config.yml
	if endpointConfig, ok := config.Endpoints[endpoint]; ok {
		if endpointConfig.AccessToken != "" {
			return endpointConfig.AccessToken
		}
	}

	// No sane default
	return ""
}

func NewClient(options ...Option) interface{} {
	c := Client{}
	for _, option := range options {
		option(&c)
	}

	// Read config eagerly, for convenience, even though we may not need it.
	// TODO: Wait until needed? Ignore error?
	config, err := readConfig()
	if err != nil {
		return errors.Wrap(err, "reading Travis config file")
	}

	// If we're given an endpoint, use it. Otherwise, use a default.
	if c.endpoint == "" {
		c.endpoint = defaultEndpoint(config)
	}

	// If we're given a token, use it.
	// Otherwise, if we're given a GitHub token, try to exchange it. Fail if we can't.
	// Otherwise, look for a default. Fail if we can't find one.
	if c.token == "" {
		if c.githubToken != "" {
			c.token, err = travisTokenFromGitHubToken(c.githubToken, c.endpoint)
			if err != nil {
				return err
			}

			c.githubToken = "" // No need to keep this around.
		} else {
			c.token = storedToken(config, c.endpoint)
		}
	}

	if c.token == "" {
		return errors.New("no access token provided")
	}

	return c
}
