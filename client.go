package travisci

import (
	"os"

	"github.com/pkg/errors"
)

var (
	orgURI = "https://api.travis-ci.org/"
	proURI = "https://api.travis-ci.com/"
)

type Client struct {
	endpoint string

	// If githubToken is set, we'll exchange it to get token.
	githubToken string
	token       string
}

// https://dave.cheney.net/2014/10/17/functional-options-for-friendly-apis
type Option func(*Client)

func AccessToken(token string) Option {
	if token == "" {
		panic(errors.Errorf("empty access token"))
	}

	return func(c *Client) {
		c.token = token
		c.githubToken = ""
	}
}

func GitHubAccessToken(token string) Option {
	if token == "" {
		panic(errors.Errorf("empty GitHub access token"))
	}

	return func(c *Client) {
		c.githubToken = token

		// Set token to a non-empty value that we'll overwrite later.
		// This way we won't try to read a token e.g. from config.
		c.token = "github"
	}
}

func Endpoint(endpoint string) Option {
	if endpoint == "" {
		panic(errors.Errorf("empty endpoint"))
	}

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
	config, err := readTravisConfig()
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
			c.token = endpointConfig.AccessToken
		}
	}

	if c.token == "" {
		return errors.New("no access token provided")
	}

	return nil
}

func NewClient(options ...Option) (*Client, error) {
	c := &Client{}
	for _, option := range options {
		option(c)
	}

	var err error
	if err = c.configureTokenAndEndpoint(); err != nil {
		return nil, err
	}

	// If we have a GitHub access token, get a Travis token from it.
	if c.githubToken != "" {
		c.token, err = travisTokenFromGitHubToken(c.githubToken, c.endpoint)
		if err != nil {
			return nil, err
		}

		c.githubToken = "" // No need to keep this around.
	}

	return c, nil
}
