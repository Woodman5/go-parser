package main

import (
	"fmt"
	"io/ioutil"
	"regexp"
)

var delimiters = [...]string {" ",
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

var regex1 = regexp.MustCompile(`(?mi)(<script|<style)([\s\w[:punct:]]*|.*)(script>|style>)`)

func main()  {

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
	text = regex1.ReplaceAllString(text, substitution)
	err = save(text)
	if err != nil {
		fmt.Print(err)
	}

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