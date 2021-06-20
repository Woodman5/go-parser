package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

func main() {

	logFile, err := os.OpenFile("parser.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer logFile.Close()

	errorLog := log.New(logFile, "ERROR ", log.Ldate|log.Ltime|log.Lshortfile)
	infoLog := log.New(logFile, "INFO ", log.Ldate|log.Ltime)

	infoLog.Print("Parser started")

	fmt.Println("Please, provide URLs one per line. Type 'exit' to finish.")

	var page htmlPage

	for fmt.Scan(&page.url); page.url != "exit"; fmt.Scan(&page.url) {

		if strings.HasPrefix(page.url, "http://") || strings.HasPrefix(page.url, "https://") {

			infoLog.Printf("Parsing URL: %s", page.url)

			res, err := http.Get(page.url)
			if err != nil {
				errorLog.Printf("GET query error: %s", err)
				fmt.Println("GET query error, possibly wrong URL")
				continue
			}
			page.setName()

			page.data, err = io.ReadAll(res.Body)
			if err != nil {
				errorLog.Printf("Response body reading error: %s", err)
				fmt.Println("Response body reading error")
				continue
			}
			err = res.Body.Close()
			if err != nil {
				errorLog.Printf("Body closing error: %s", err)
				fmt.Println("Body closing error")
				continue
			}

			err = page.save()
			if err != nil {
				errorLog.Printf("Error saving HTML file: %s", err)
				fmt.Println("Error saving HTML file")
			}

			page.makeMap()

			err = page.saveMap()
			if err != nil {
				errorLog.Printf("Error saving MAP file: %s", err)
				fmt.Println("Error saving MAP file")
			}

			err = page.saveToDB()
			if err != nil {
				errorLog.Printf("Error saving data to database: %s", err)
				fmt.Println("Error saving data to database")
			}

			fmt.Printf("Parsing done. Found %d different words.\n", len(page.dict))

		} else {
			infoLog.Printf("Entered wrong URL - %s", page.url)
			fmt.Println("Url must start with 'http://' or 'https://'")
		}

	}
	infoLog.Print("Parser stopped")
}
