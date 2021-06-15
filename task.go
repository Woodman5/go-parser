package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
	// "strings"
)

var delimiters = [...]string{" ",
	",",
	".",
	"! ",
	"?",
	"\"",
	";",
	":",
	"[",
	"]",
	"(",
	")",
	"\n",
	"\r",
	"\t",
}
var substitution = ""

var regex1 = regexp.MustCompile(`(?mi)(<script[>\s]*|<style)([\s\w[:punct:]]*|.*)([<\s]*/script>|style>)`)
var regex2 = regexp.MustCompile(`(?m)(<!-- )|( -->)`)
var regex3 = regexp.MustCompile(`(?m)<(/?[^>]+)>`)
var regex4 = regexp.MustCompile(`(?m)(description"\s+content=")|(title" content=")`)
var regex5 = regexp.MustCompile(`(?m)(\s\-\s)|([\[\]])`)
var regex6 = regexp.MustCompile(`(?m)([\s—«»\(\)\?!/>"\{\}\.,:;']+)|(&#160;)`)

func main() {

	var address string

	fmt.Println("Please, provide URLs one per line. Type 'exit' to finish.")

	for fmt.Scan(&address); address != "exit"; fmt.Scan(&address) {
		if strings.HasPrefix(address, "http://") || strings.HasPrefix(address, "https://") {
			res, err := http.Get(address)
			if err != nil {
				fmt.Print(err)
			}
			defer res.Body.Close()
			body, err := io.ReadAll(res.Body)
			if err != nil {
				fmt.Print(err)
			}

			processAdress(string(body))

		} else {
			fmt.Println("Url must start with 'http://' or 'https://'")
		}

	}
}

// func loadPage() (string, error) {
// 	body, err := ioutil.ReadFile("text.html")
// 	if err != nil {
// 		return "", err
// 	}
// 	return string(body), nil
// }

func save(text string, dict map[string]int) error {
	err := ioutil.WriteFile("text.txt", []byte(text), 0600)
	b := new(bytes.Buffer)
	for key, value := range dict {
		fmt.Fprintf(b, "%s - %d\n", key, value)
	}
	err = ioutil.WriteFile("dict.txt", []byte(b.String()), 0600)
	return err
}

func processAdress(text string) {
	text = regex4.ReplaceAllString(text, "/>")
	text = regex1.ReplaceAllString(text, substitution)
	text = regex2.ReplaceAllString(text, substitution)
	text = regex3.ReplaceAllString(text, " ")
	text = regex5.ReplaceAllString(text, " ")
	text = strings.Replace(text, "&#37;", "%", -1)
	text = strings.Replace(text, "&#8209;", "-", -1)
	text = strings.Replace(text, "&quot;", "\"", -1)
	text = strings.Replace(text, "&nbsp;", " ", -1)
	text = regex6.ReplaceAllString(text, "\n")
	text = strings.TrimSpace(text)
	text = strings.ToLower(text)

	wordList := strings.Split(text, "\n")

	fmt.Println(len(wordList))

	m := make(map[string]int)

	for _, v := range wordList {
		m[v]++
	}

	fmt.Println(m)

	err := save(text, m)
	if err != nil {
		fmt.Print(err)
	}
}
