package auth

import (
	"errors"
	"net/http"
	"strings"
)

// Authorization: ApiKey <key>
func GrabAPIKey(headers http.Header) (string, error) {
	val := headers.Get("Authorization")
	if val == "" {
		return "", errors.New("no authentication info")
	}

	vals := strings.Split(val, " ")
	if len(vals) != 2 {
		return "", errors.New("malformed auth header")
	}

	if vals[0] != "ApiKey" {
		return "", errors.New("malformed first auth header")
	}

	return vals[1], nil
}
