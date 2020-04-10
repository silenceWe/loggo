package main

import (
	"os"
	"sync/atomic"
	"time"

	"github.com/silenceWe/loggo"
)

var log1 *loggo.Logger

func main() {
	loggo.InitDefaultLog()

	log1 = &loggo.Logger{Level: loggo.ALL}
	log1.SetWriter(&loggo.FileWriter{FileName: "./loggo1.log", Compress: true})
	start()
}

func start() {

	go ticker()
	print()

}

func print() {
	for {
		log1.Infoln("aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa")
		atomic.AddInt64(&count, 1)
	}
}

var count int64
var total int64

func ticker() {
	tk := time.NewTicker(1 * time.Second)
	for {
		select {
		case <-tk.C:
			loggo.Infofn("QPS:[%d]", count)
			atomic.AddInt64(&total, count)
			if total > 1000000 {
				os.Exit(0)
			}
			count = 0
		}
	}
}

/**
len : 55
[2020-04-10 14:46:59.975|INFO ] QPS:[174230]
[2020-04-10 14:47:00.974|INFO ] QPS:[171390]
[2020-04-10 14:47:01.974|INFO ] QPS:[173086]
[2020-04-10 14:47:02.974|INFO ] QPS:[150851]
[2020-04-10 14:47:03.974|INFO ] QPS:[169345]
[2020-04-10 14:47:04.974|INFO ] QPS:[171737]
**/
