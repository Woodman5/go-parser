package main

import (
	"fmt"
	"io/ioutil"
	"regexp"
	"strings"
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
var regex7 = regexp.MustCompile(`(?m)([\s—«»\(\)\?!/>"\.,:;]+)|(&#160;)`)

func main() {

	//res, err := http.Get("https://www.simbirsoft.com")
	//if err != nil {
	//	fmt.Print(err)
	//}
	//defer res.Body.Close()
	//body, err := io.ReadAll(res.Body)
	text, err := loadPage()
	if err != nil {
		fmt.Print(err)
	}
	//fmt.Print(string(text))
	text = regex4.ReplaceAllString(text, "/>")
	text = regex1.ReplaceAllString(text, substitution)
	text = regex2.ReplaceAllString(text, substitution)
	text = regex3.ReplaceAllString(text, " ")
	text = regex5.ReplaceAllString(text, " ")
	text = strings.Replace(text, "&#37;", "%", -1)
	text = strings.Replace(text, "&#8209;", "-", -1)
	text = regex7.ReplaceAllString(text, "\n")
	text = strings.TrimSpace(text)
	err = save(text)
	if err != nil {
		fmt.Print(err)
	}
	test := [4]int{1, 2, 3, 4}
	wordList := strings.Split(text, "\n")
	fmt.Printf("%T\n", test)
	fmt.Printf("%T\n", wordList)
	fmt.Println(wordList)
	fmt.Println(len(wordList))
}

func loadPage() (string, error) {
	body, err := ioutil.ReadFile("text.html")
	if err != nil {
		return "", err
	}
	return string(body), nil
}

func save(text string) error {
	return ioutil.WriteFile("text.txt", []byte(text), 0600)
}
