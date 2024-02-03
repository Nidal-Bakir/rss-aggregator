package utils

import "net/url"

func IsValidUrl(str string) bool {
	u, err := url.ParseRequestURI(str)
	return err == nil && u.Scheme != "" && u.Host != ""
}

func Clamp(x int, min int, max int) int {
	if x <= min {
		return min
	}
	if x >= max {
		return max
	}
	return x
}
