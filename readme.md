`kh` is a programmatic way to check for the validity of API tokens or webhooks. The services against
which it is able to check originally came from the popular [keyhack](https://github.com/streaak/keyhacks#Slack-API-token)
repo by [@streaak](https://github.com/streaak/).

## Usage

```bash
$ kh github-token XXXXXXXXXXXXXXXXXXXXXXXXX

$ ./my-custom-token-scanner | kh slack-token - | tee -a valid_slack_tokens.txt

$ xargs kh slack-token < maybe_tokens.txt| tee -a valid_slack_tokens.txt
```

If the token is valid, `kh` will print the token and return a 0 status to bash. If the token is
invalid, nothing will be printed and the status returned will be 1. The output is minimal so that
the tool can be used in existing workflows, bash pipelines and scripts.

## Expandability

It's possible to add services to the tool by modifying the configuration YAML file. 

```yaml
# Demo Service With All Params
sass-api:
  name: sass-api
  request:
    method: POST # [REQUIRED]
    url: 'https://sass-api.io/api/auth' # [REQUIRED]
    headers:
      Authorization: Bearer %s
  validator: # [REQUIRED if 200/40x http status is not indicative of success/failure]
    custom: true
```

In the parameters where a token is to be interpolated, place a template symbol, `%s`, in place of
the token value.

By default, `kh` will declare a token as valid if the API returns a 200 HTTP status. Not all APIs are
create equal nor do they use semantic HTTP status codes when replying. If you're attempting to add a
new service to `kh` and both valid and invalid tokens return a `200`, then a custom validator must be written.

In addition to editing the configuration YAML, users must add the subcommand to the `/cmd`
folder in this repository's root. When declaring a custom validator in the YAML file, users must also 
define what a valid response looks like

```go
// each subcommand's init function must add the subcommand to the root cli command
// and then add the validator function to the keyhack registry so that it knows
// what a good http response looks like
func init() {
	rootCmd.AddCommand(slackTokenCmd)
	keyhack.Registry["slack-token"].Validator.Fn = validateSlack
}

// ensure the command name matches the entry in the YAML file
var slackTokenCmd = newCommand("slack-token", "Checks a token against the Slack API")

// validator functions define what a successful authentication means 
// based on the http response of the API call issued by keyhacks
func validateSlack(resp *http.Response) (ok bool, err error) {
	ok = resp.Header["X-Oauth-Scopes"] != nil
	return
}
```

If you don't need a custom validator, that is, if the API returns anything but a 200 with invalid creds, then the following is all that's needed in the new service:

```go
// cmd/github.go
package cli

func init() {
	githubTokenCmd := newCommand("github-token", "Checks a token against the GitHub API")
	rootCmd.AddCommand(githubTokenCmd)
}
```


## Structure

```
├── cmd			# this is where new plugins go
│   ├── cli.go		# main entry point logic for the CLI utility
│   └── <more services here>
├── go.mod
├── go.sum
├── keyhacks.yml	# tool configuration; add new service definitions here
├── main.go		
├── pkg
│   └── keyhack
│       └── keyhack.go	# core keyhack framework logic

```
