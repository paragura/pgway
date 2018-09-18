package util

func UrlSanitize(pathStr string) string {
	var retRunes []rune
	var isSlash = false
	for _, pathRune := range pathStr {
		if pathRune == '/' {
			if isSlash {
				continue
			} else {
				isSlash = true
			}
		} else {
			isSlash = false
		}
		retRunes = append(retRunes, pathRune)
	}
	return string(retRunes)
}
