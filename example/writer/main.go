package main

import (
	"github.com/silenceWe/loggo/writer"
	"log"
)

func main(){
	wtr1 := writer.NewLogWriter(&writer.WriterOption{StdOut: true,  MaxAge: 7})
	wtr2 := writer.NewLogWriter(&writer.WriterOption{StdOut: true,  MaxAge: 7,FileName:"./log2.log"})
	logger1 := log.Logger{}
	logger1.SetOutput(wtr1)
	logger1.SetFlags(log.LUTC)
	logger1.Println("test1")

	logger2 := log.Logger{}
	logger2.SetOutput(wtr2)
	logger2.SetFlags(log.LUTC)
	logger2.Println("test2")
}
