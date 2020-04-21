package main

import (
	"os"

	"github.com/neoul/gostudy/datablock/log"
)

func main() {
	l := log.NewLog("hello", os.Stdout)
	l.Debug("DEBUG message")
}