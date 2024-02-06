package main

import (
	"fmt"
	"log"
	"os"
	"time"
)

var (
	WarningLogger *log.Logger
	InfoLogger    *log.Logger
	ErrorLogger   *log.Logger
)

func init() {

	rightnow := time.Now().In(time.Local)
	outtaTime := rightnow.Format("2006-01-02")

	filename := "tagger-" + outtaTime + ".log"

	file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		fmt.Println("Failed to initiate log file: " + filename)
		log.Fatal(err)
	}

	InfoLogger = log.New(file, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	WarningLogger = log.New(file, "WARNING: ", log.Ldate|log.Ltime|log.Lshortfile)
	ErrorLogger = log.New(file, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
}

func Logit(msg string, debug bool, state string) {

	switch state {
	case "warn":
		WarningLogger.Println(msg)
	case "info":
		InfoLogger.Println(msg)
	case "err":
		ErrorLogger.Println(msg)
	}

	if debug {
		fmt.Println(msg)
	}
}
