package main

import (
	"github.com/silenceWe/loggo"
	"log"
)

func main() {
	writer := &loggo.FileWriter{MaxAge: 7, FileName: "./log2.log"}
	writer.Init()
	logger1 := log.Logger{}
	logger1.SetOutput(writer)
	logger1.SetFlags(log.LUTC)
	logger1.Println("test1")

	loggo1 := &loggo.Logger{Level: loggo.ALL}
	loggo1.SetWriter(&loggo.FileWriter{FileName: "./loggo1.log"})
	loggo1.Infoln("test")
}
