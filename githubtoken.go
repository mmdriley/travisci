package travisci

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/pkg/errors"
)

type githubAuthRequest struct {
	GitHubToken string `json:"github_token"`
}

type githubAuthResponse struct {
	AccessToken string `json:"access_token"`
}

func travisTokenFromGitHubToken(githubAccessToken string, endpoint string) (string, error) {
	requestBody := githubAuthRequest{GitHubToken: githubAccessToken}

	url, err := url.Parse(endpoint)
	if err != nil {
		return "", errors.Wrapf(err, "parsing endpoint URL %v", endpoint)
	}
	url.Path = "/auth/github"

	request, err := http.NewRequest("POST", url.String(), bytes.NewReader(mustJSONMarshal(requestBody)))
	if err != nil {
		return "", errors.Wrap(err, "creating request")
	}

	request.Header.Set("Content-Type", "application/json")

	// Use a User-Agent that begins with "Travis" to avoid https://github.com/travis-ci/travis-ci/issues/5649
	// This only seems to a be a problem at travis-ci.org (not .com)
	request.Header.Set("User-Agent", "Travis")

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return "", errors.Wrapf(err, "POST to %v", url.String())
	}
	defer response.Body.Close()

	if response.StatusCode != 200 {
		// Return full response body in error to allow for debugging.
		responseBytes, err := ioutil.ReadAll(response.Body)
		if err != nil {
			responseBytes = nil
		}

		return "", errors.Errorf(`"%s" from GitHub auth endpoint: %s`, response.Status, responseBytes)
	}

	var responseBody githubAuthResponse
	if err := json.NewDecoder(response.Body).Decode(&responseBody); err != nil {
		return "", errors.Wrapf(err, "parsing GitHub auth response as JSON")
	}

	return responseBody.AccessToken, nil
}
