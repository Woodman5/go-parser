package main

import (
	"fmt"
	"io"
	"net/http"
	"strings"
)

func main() {
	initLogger()
	log.Info("Parser started")

	fmt.Println("Please, provide URLs one per line. Type 'exit' to finish.")

	var url string

	for fmt.Scan(&url); url != "exit"; fmt.Scan(&url) {

		url = strings.ToLower(url)

		if strings.HasPrefix(url, "http://") || strings.HasPrefix(url, "https://") {

			log.Info("Parsing URL: ", url)

			res, err := http.Get(url)
			if err != nil {
				log.Error("GET query error: ", err)
				fmt.Println("GET query error, possibly wrong URL")
				continue
			}

			data, err := io.ReadAll(res.Body)
			if err != nil {
				log.Error("Response body reading error: ", err)
				fmt.Println("Response body reading error")
				continue
			}

			err = res.Body.Close()
			if err != nil {
				log.Error("Body closing error: ", err)
				fmt.Println("Body closing error")
				continue
			}

			page := newPage(url, data)

			err = page.save()
			if err != nil {
				log.Error("Error saving HTML file: ", err)
				fmt.Println("Error saving HTML file")
			}

			err = page.saveMap()
			if err != nil {
				log.Error("Error saving MAP file: ", err)
				fmt.Println("Error saving MAP file")
			}

			err = page.saveToDB()
			if err != nil {
				log.Error("Error saving data to database: %s", err)
				fmt.Println("Error saving data to database")
			}

			msg := fmt.Sprintf("Parsing done. Found %d different words.", len(page.dict))
			fmt.Println(msg)
			log.Info(msg)

		} else {
			log.Error("Entered wrong URL - ", url)
			fmt.Println("Url must start with 'http://' or 'https://'")
		}

	}
	log.Info("Parser stopped")
}
