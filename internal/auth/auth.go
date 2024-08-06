package auth

import (
	"errors"
	"net/http"
	"strings"
)

func ParseAuthHeader(headers http.Header, content string) (string, error) {
	authHeader := headers.Get("Authorization")
	if authHeader == "" {
		return "", errors.New("no auth header included in request")
	}

	splitAuth := strings.Split(authHeader, " ")
	if len(splitAuth) < 2 || splitAuth[0] != content {
		return "", errors.New("malformed authorization header")
	}

	return splitAuth[1], nil
}
