package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"io/ioutil"
	"regexp"
	"strings"
	"time"
)

var regex1 = regexp.MustCompile(`(?mi)(<script[>\s]*|<style)([\s\w[:punct:]]*|.*)([<\s]*/script>|style>)`)
var regex2 = regexp.MustCompile(`(?m)(<!-- )|( -->)`)
var regex3 = regexp.MustCompile(`(?m)<(/?[^>]+)>`)
var regex4 = regexp.MustCompile(`(?m)(description"\s+content=")|(title" content=")`)
var regex5 = regexp.MustCompile(`(?m)(\s\-\s)|([\[\]])`)
var regex6 = regexp.MustCompile(`(?m)([\s—«»\\()?!/>"{},:;']+)|(&#160;)|(\.\s)`)
var regex7 = regexp.MustCompile(`(?mi)(http://)|(https://)`)
var regex8 = regexp.MustCompile(`(?m)\n{2,}`)

type htmlPage struct {
	url, fileName string
	data          []byte
	dict          map[string]int
	date          int64
}

func (p htmlPage) save() error {
	err := ioutil.WriteFile(p.fileName+".html", p.data, 0600)
	return err
}

func (p *htmlPage) setName() {
	p.date = time.Now().Unix()
	name := regex7.ReplaceAllString(p.url, "")
	name = strings.Replace(name, ".", "-", -1)
	index := strings.Index(name, "/")
	if index != -1 {
		name = name[:index]
	}
	p.fileName = fmt.Sprintf("%s-%d", name, p.date)
}

func (p htmlPage) saveMap() error {
	b := new(bytes.Buffer)
	for key, value := range p.dict {
		fmt.Fprintf(b, "%s - %d\n", key, value)
	}
	err := ioutil.WriteFile(p.fileName+".txt", []byte(b.String()), 0600)
	return err
}

func (p htmlPage) convertToJSON() ([]byte, error) {
	jsonData, err := json.Marshal(p.dict)
	return jsonData, err
}

func (p htmlPage) saveToDB() error {
	db, err := sql.Open("sqlite3", "words.sqlite")
	if err != nil {
		return err
	}
	defer db.Close()

	jsonData, err := p.convertToJSON()
	if err != nil {
		return err
	}
	_, err = db.Exec("insert into words (url, filename, words, date) values ($1, $2, $3, $4)",
		p.url, p.fileName, string(jsonData), p.date)

	return err
}

func (p *htmlPage) makeMap() {
	text := string(p.data)
	text = regex4.ReplaceAllString(text, "/>")
	text = regex1.ReplaceAllString(text, "")
	text = regex2.ReplaceAllString(text, "")
	text = regex3.ReplaceAllString(text, " ")
	text = regex5.ReplaceAllString(text, " ")
	text = strings.Replace(text, "&#37;", "%", -1)
	text = strings.Replace(text, "&#8209;", "-", -1)
	text = strings.Replace(text, "&quot;", "\"", -1)
	text = strings.Replace(text, "&nbsp;", " ", -1)
	text = regex6.ReplaceAllString(text, "\n")
	text = regex8.ReplaceAllString(text, "\n")
	text = strings.TrimSpace(text)
	text = strings.ToLower(text)

	wordList := strings.Split(text, "\n")

	p.dict = make(map[string]int)

	for _, v := range wordList {
		p.dict[v]++
	}
}
