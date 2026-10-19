package auth

import (
	"errors"
	"net/http"
	"strings"
)

func GetAPIKey(headers http.Header) (string, error) {
	authHeader := headers.Get("Authorization")
	if authHeader == "" {
		return "", errors.New("Authorization header is missing")
	}

	hasPrefix := strings.HasPrefix(authHeader, "ApiKey")
	if !hasPrefix {
		return "", errors.New("Invalid Authorization header format")
	}

	apiString := strings.TrimPrefix(authHeader, "ApiKey")
	apiString = strings.TrimSpace(apiString)

	if apiString == "" {
		return "", errors.New("Invalid API key")
	}

	return apiString, nil
}
