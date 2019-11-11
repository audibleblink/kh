package keyhack

import (
	"net/http"
)

func validateGithub(resp *http.Response) (ok bool, err error) {
	ok = resp.StatusCode == 200
	return
}

