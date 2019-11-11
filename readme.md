`kh` is a programatic way to check for the validity of API tokens or webhooks. The services against
which it is able to check originally came from the popular [keyhack](https://github.com/streaak/keyhacks#Slack-API-token)
repo by [@streaak](https://github.com/streaak/).

## Usage

```bash
$ kh github-token XXXXXXXXXXXXXXXXXXXXXXXXX
=> XXXXXXXXXXXXXXXXXXXXXXXXX
```

If the token is valid, `kh` will print the token and return a 0 status to bash. If the token is
invalid, nothing will be printed and the status returned will be 1. The output is minimal so that
the tool can be used in existing workflows, bash pipelines and scripts.

## Expandability

It's possible to add services to the tool by modifying the configuration YAML file. 

```yaml
github-token:
  name: github-token
  request:
    method: GET
    url: 'https://api.github.com/users'
    headers:
      Authorization: "token %s"

slack-token:
  name: slack-token
  request:
    method: POST
    url: 'https://slack.com/api/auth.test?token=%s&pretty=1'
```

In the parameters where a token is to be interpolated, place a template symbol, `%s`, in place of
the token value.

In addition to the edits to the configuration YAML, users must add the subcommand to the `/cmd`
folder in this repository's root. Users must also define what a validated response means. 

```go
// each subcommand's init function must add the subcommand to the root cli command
// and then add the validator function to the keyhack registry so that it knows
// what a good http response looks like
func init() {
	rootCmd.AddCommand(slackTokenCmd)
	keyhack.Registry["slack-token"].Validator = validateSlack
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
