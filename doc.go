package travisci

/*

write token to ~/.travis/config.yml with:
TRAVIS_TOKEN=abc travis whoami --pro

default endpoint:
travis endpoint --set-default --pro

What scopes are needed for GitHub authentication?
curl https://api.travis-ci.com/v3/ | jq .config.github.scopes

*/
