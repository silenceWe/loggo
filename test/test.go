package main

import (
	"os"
	"sync/atomic"
	"time"

	"github.com/silenceWe/loggo"
)

var log1 *loggo.Logger

func main() {

	loggo.InitDefaultLog(&loggo.LoggerOption{StdOut: true, Level: loggo.ALL})
	log1 = loggo.NewLoggo(&loggo.LoggerOption{StdOut: false, Level: loggo.ALL, FileName: "./log/log1.log", MaxSize: 500})
	start()
}

func start() {

	go ticker()
	print()

}

func print() {
	for {
		log1.Infoln("aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa", "bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb", "ccccccccccccccccccccccccccccccccccccccccccccccc", "ddddddddddddddddddddddddddddddddddddddddddd")
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
len : 4
[2020-01-20 17:06:54.391|INFO ] QPS:[132076]
[2020-01-20 17:06:55.390|INFO ] QPS:[135453]
[2020-01-20 17:06:56.390|INFO ] QPS:[125967]
[2020-01-20 17:06:57.391|INFO ] QPS:[134366]
[2020-01-20 17:06:58.390|INFO ] QPS:[134460]
[2020-01-20 17:06:59.390|INFO ] QPS:[131881]
[2020-01-20 17:07:00.391|INFO ] QPS:[132926]
[2020-01-20 17:07:01.390|INFO ] QPS:[132172]


len : 200
[2020-01-20 17:10:59.456|INFO ] QPS:[103401]
[2020-01-20 17:11:00.456|INFO ] QPS:[97803]
[2020-01-20 17:11:01.456|INFO ] QPS:[97878]
[2020-01-20 17:11:02.456|INFO ] QPS:[102290]
[2020-01-20 17:11:03.456|INFO ] QPS:[102914]
[2020-01-20 17:11:04.456|INFO ] QPS:[107988]
[2020-01-20 17:11:05.456|INFO ] QPS:[106840]
[2020-01-20 17:11:06.456|INFO ] QPS:[114852]
[2020-01-20 17:11:07.456|INFO ] QPS:[108525]
[2020-01-20 17:11:08.456|INFO ] QPS:[103544]

len : 200 muil str
[2020-01-20 17:28:04.832|INFO ] QPS:[127633]
[2020-01-20 17:28:05.832|INFO ] QPS:[136644]
[2020-01-20 17:28:06.832|INFO ] QPS:[132268]
[2020-01-20 17:28:07.832|INFO ] QPS:[132977]
[2020-01-20 17:28:08.832|INFO ] QPS:[132456]
[2020-01-20 17:28:09.832|INFO ] QPS:[128749]
[2020-01-20 17:28:10.832|INFO ] QPS:[131483]
[2020-01-20 17:28:11.832|INFO ] QPS:[121070]
**/
