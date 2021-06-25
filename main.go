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

	var page htmlPage

	for fmt.Scan(&page.url); page.url != "exit"; fmt.Scan(&page.url) {

		page.url = strings.ToLower(page.url)

		if strings.HasPrefix(page.url, "http://") || strings.HasPrefix(page.url, "https://") {

			log.Info("Parsing URL: ", page.url)

			res, err := http.Get(page.url)
			if err != nil {
				log.Error("GET query error: ", err)
				fmt.Println("GET query error, possibly wrong URL")
				continue
			}
			page.setName()

			page.data, err = io.ReadAll(res.Body)
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

			err = page.save()
			if err != nil {
				log.Error("Error saving HTML file: ", err)
				fmt.Println("Error saving HTML file")
			}

			page.makeMap()

			err = page.saveMap()
			if err != nil {
				log.Error("Error saving MAP file: ", err)
				fmt.Println("Error saving MAP file")
			}

			// err = page.saveToDB()
			// if err != nil {
			// 	log.Error("Error saving data to database: %s", err)
			// 	fmt.Println("Error saving data to database")
			// }

			fmt.Printf("Parsing done. Found %d different words.\n", len(page.dict))

		} else {
			log.Error("Entered wrong URL - ", page.url)
			fmt.Println("Url must start with 'http://' or 'https://'")
		}

	}
	log.Info("Parser stopped")
}
