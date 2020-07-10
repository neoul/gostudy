package main

import (
	"fmt"

	"github.com/neoul/gostudy/samples"
)

func main() {
	samples.Switch1()
	samples.Switch2()
	fmt.Println(samples.Switch3('A'))
	samples.Switch4_goto()

	samples.Timer1()
	samples.Ticker1()
	samples.WorkerPools()
	samples.Select1()

}
