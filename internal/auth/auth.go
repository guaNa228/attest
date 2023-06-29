package auth

import (
	"errors"
	"net/http"
	"strings"
)

func GetJWT(headers http.Header) (string, error) {
	val := headers.Get("Authorization")
	if val == "" {
		return "", errors.New("no authentication info found")
	}

	vals := strings.Split(val, " ")
	if len(vals) != 2 {
		return "", errors.New("malformed auth header")
	}

	if vals[0] != "JWT" {
		return "", errors.New("malformed auth header key")
	}

	return vals[1], nil
}
