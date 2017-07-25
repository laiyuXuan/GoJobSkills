package utils

import "regexp"

func RemoveBlanks(body string) string {
	emptyRx := "\\s+"
	emptyCompile := regexp.MustCompile(emptyRx)
	body = emptyCompile.ReplaceAllString(body, "")
	return body
}

func RemoveSpace(matched string) string {
	htmlSpaceRx := "&nbsp"
	htmlSpaceComplie := regexp.MustCompile(htmlSpaceRx)
	matched = htmlSpaceComplie.ReplaceAllString(matched, "")
	return matched
}

func RemoveHtmlTag(matched string) string {
	htmlLabelRx := "<.+?>"
	htmlCompile := regexp.MustCompile(htmlLabelRx)
	matched = htmlCompile.ReplaceAllString(matched, "")
	return matched
}
