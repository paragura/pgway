package util

import "regexp"

func UrlSanitize(urlStr string) string {
	reg := regexp.MustCompile(`/+`)
	return reg.ReplaceAllString(urlStr, "/")
}
