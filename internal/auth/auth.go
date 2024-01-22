package auth

import (
	"errors"
	"net/http"
	"strings"
)

// Authorization: apiKey <some_api_key>
func GetApiKey(header http.Header) (apiKey string, err error) {
	authorization := header.Get("Authorization")
	if authorization == "" {
		return "", errors.New("Authorization header not provided")
	}

	result := strings.Split(authorization, " ")
	if len(result) != 2 {
		return "", errors.New("malformed Authorization header")

	}

}
