package main

import (
	"fmt"

	"github.com/neoul/gostudy/statements"
)

func main() {
	statements.Switch1()
	statements.Switch2()
	fmt.Println(statements.Switch3('A'))
	statements.Switch4_goto()

	statements.Timer1()
	statements.Ticker1()
	statements.WorkerPools()
	statements.Select1()
}
