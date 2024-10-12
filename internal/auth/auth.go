package auth

import (
	"errors"
	"net/http"
	"strings"
)

// Extracts API key from the headers
func GetApiKey(headers http.Header) (string, error) {
	val := headers.Get("Authorization")
	if val == "" {
		return "", errors.New("missing API Key")
	}

	vals := strings.Split(val, " ")
	if len(vals) != 2 {
		return "", errors.New("missing API Key")
	}

	if vals[0] != "ApiKey" {
		return "", errors.New("malformed first part of Auth header")
	}

	return vals[1], nil
}
