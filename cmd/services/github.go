package services

import (
	cli "github.com/audibleblink/kh/cmd"
)

// Service configuration
const (
	GithubTokenSubCmd = "github-token"
	GithubTokenToken  = "<token>"

	GithubOauthSubCmd = "github-oauth"
	GithubOauthToken  = "<client_id:client_secret>"
)

// init registers the GitHub services
func init() {
	// Create commands for both GitHub services
	_ = cli.NewServiceCommand(GithubTokenSubCmd, GithubTokenToken)
	_ = cli.NewServiceCommand(GithubOauthSubCmd, GithubOauthToken)

	// No custom validators needed - both use default 200 OK check
}
