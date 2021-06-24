package main

import (
	"github.com/sirupsen/logrus"
	"os"
	"strings"
)

type Formatter struct{}

func (f *Formatter) Format(entry *logrus.Entry) ([]byte, error) {
	output := "[%lvl%]: %time% - %msg%\n"

	timestampFormat := "2006-01-02 15:04:05 Z07:00 MST"

	output = strings.Replace(output, "%time%", entry.Time.Format(timestampFormat), 1)
	output = strings.Replace(output, "%msg%", entry.Message, 1)
	level := strings.ToUpper(entry.Level.String())
	output = strings.Replace(output, "%lvl%", level, 1)

	return []byte(output), nil
}

var log = logrus.New()

func initLogger() {
	log.SetFormatter(&Formatter{})
	log.SetLevel(logrus.DebugLevel)
	log.SetOutput(os.Stdout)
	file, err := os.OpenFile("parser.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err == nil {
		log.SetOutput(file)
	} else {
		log.Info("Failed to log to file, using default stderr")
	}
}
