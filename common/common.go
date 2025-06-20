package common

import (
	"log"
	"runtime"
)

type Request map[string]interface{}

type Response struct {
	Data  interface{}
	Error error
}

const (
	KEY      = "KEY"
	PASSWORD = "PASSWORD"
)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
}
