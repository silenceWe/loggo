package main

import "loggo/loggo"

type logConfig struct {
	DefaultLogConfig *loggo.LoggerOption
}

var LogConfig = &logConfig{}

func init() {
}
