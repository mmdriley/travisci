package travisci

import (
	"io/ioutil"
	"os"
	"os/user"
	"path"

	"github.com/go-yaml/yaml"
	"github.com/pkg/errors"
)

// Structure of ~/.travis/config.yml
type config struct {
	Endpoints       map[string]configEndpoint
	DefaultEndpoint string `yaml:"default_endpoint"`
}

type configEndpoint struct {
	AccessToken string `yaml:"access_token"`
}

func configPath() (string, error) {
	if path := os.Getenv("TRAVIS_CONFIG_PATH"); path != "" {
		return path, nil
	}

	currentUser, err := user.Current()
	if err != nil {
		return "", err
	}

	return path.Join(currentUser.HomeDir, ".travis"), nil
}

func readTravisConfig() (config, error) {
	configPath, err := configPath()
	if err != nil {
		return config{}, err
	}

	name := path.Join(configPath, "config.yml")

	f, err := os.Open(name)
	if os.IsNotExist(err) {
		return config{}, nil
	} else if err != nil {
		return config{}, errors.Wrapf(err, "opening %v", name)
	}

	bytes, err := ioutil.ReadAll(f)
	if err != nil {
		return config{}, errors.Wrapf(err, "reading %v", name)
	}

	var result config
	err = yaml.Unmarshal(bytes, &result)
	if err != nil {
		return config{}, errors.Wrapf(err, "parsing %v as yaml", name)
	}

	return result, nil
}
