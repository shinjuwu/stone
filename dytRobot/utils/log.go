package utils

import (
	"log"
	"os"
	"time"
)

var setlevel = LOG_INFO

const (
	LOG_DEBUG int = iota
	LOG_INFO
	LOG_ERROR
)

func LogInfo(logLevel int, message string, channel chan string) {
	if logLevel < setlevel {
		return
	}

	content := time.Now().Format("2006/01/02 15:04:05")
	switch logLevel {
	case LOG_DEBUG:
		content += " [DEBUG] "
	case LOG_INFO:
		content += " [INFO] "
	case LOG_ERROR:
		content += " [ERROR] "
	}
	content += message + "\n"
	WriteToFile(content)
	if channel != nil {
		channel <- content
	}
}

func WriteToFile(content string) {
	_, err := os.Stat("./log")
	if err != nil && os.IsNotExist(err) {
		os.Mkdir("./log", os.ModePerm)
	}

	date := time.Now().Format("2006_01_02")
	filePath := "./log/" + date + ".log"
	data := []byte(content)
	option := os.O_WRONLY | os.O_APPEND | os.O_CREATE
	file, err := os.OpenFile(filePath, option, 0666)
	if err != nil {
		log.Fatal(err)
	}
	_, err = file.Write(data)
	if err != nil {
		log.Fatal(err)
	}
}
