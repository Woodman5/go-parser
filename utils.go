package main

import (
	"regexp"
	"strings"
)

var htmlCodes = [...]string{"&#37;", "%", "&#8209;", "-", "&quot;", "\"", "&nbsp;", " "}

var regex = [...]string{
	`(?m)(description"\s+content=")|(title" content=")`, "/>",
	`(?mi)(<script[>\s]*|<style)([\s\w[:punct:]]*|.*)([<\s]*/script>|style>)`, "",
	`(?m)(<!-- )|( -->)`, "",
	`(?m)<(/?[^>]+)>`, " ",
	`(?m)(\s\-\s)|([\[\]])`, " ",
	`(?m)([\s—«»\\()?!/>"{},:;']+)|(&#160;)|(\.\s)`, "\n",
	`(?m)\n{2,}`, "\n",
}

func cleanHtmlFromPage(data []byte) string {
	text := string(data)

	for i := 0; i < len(htmlCodes)-1; i = i + 2 {
		text = strings.Replace(text, htmlCodes[i], htmlCodes[i+1], -1)
	}

	for i := 0; i < len(regex)-1; i = i + 2 {
		text = regexp.MustCompile(regex[i]).ReplaceAllString(text, regex[i+1])
	}

	text = strings.TrimSpace(text)
	text = strings.ToLower(text)
	return text
}
