package main

import (
	"os"

	"github.com/neoul/gostudy/datablock/log"
	yref "github.com/neoul/gostudy/datablock/reflection"
)


func main() {
	// var i int
	var ri int
	a := 1
	ylog := log.NewLog("main", os.Stdout)
	ri = yref.NewValue(a, 10).(int)
	ylog.Debug(yref.PrintValue(ri))
}