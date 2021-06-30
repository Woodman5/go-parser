package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"io/ioutil"
	"strings"
	"time"
)

type htmlPage struct {
	url, fileName string
	data          []byte
	dict          map[string]int
	date          int64
}

func newPage(url string, data []byte) *htmlPage {
	p := new(htmlPage)
	p.url = url
	p.data = data
	p.setName()
	p.makeMap()
	return p
}

func (p htmlPage) save() error {
	err := ioutil.WriteFile("pages/"+p.fileName+".html", p.data, 0600)
	return err
}

func (p *htmlPage) setName() {
	p.date = time.Now().Unix()

	name := strings.TrimPrefix(p.url, "http://")
	name = strings.TrimPrefix(name, "https://")
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
	err := ioutil.WriteFile("pages/"+p.fileName+".txt", b.Bytes(), 0600)
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
	text := cleanHtmlFromPage(p.data)

	wordList := strings.Split(text, "\n")

	p.dict = make(map[string]int)

	for _, v := range wordList {
		p.dict[v]++
	}
}
