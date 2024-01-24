package auth

import (
	"errors"
	"net/http"
	"strings"
)

// Authorization: ApiKey <some_api_key>
func GetApiKey(header http.Header) (apiKey string, err error) {
	authorization := header.Get("Authorization")
	if authorization == "" {
		return "", errors.New("authorization header not provided")
	}

	result := strings.Split(authorization, " ")
	if len(result) != 2 || result[0] != "ApiKey" || len(result[1]) != 64 {
		return "", errors.New("malformed Authorization header")
	}

	return result[1], nil

}
