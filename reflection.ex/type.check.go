package main

import (
	"log"
	"strconv"
)

func Sprint(x interface{}) string {
	type stringer interface {
		String() string
	}
	switch x := x.(type) {
	case stringer:
		return x.String()
	case string:
		return x
	case int:
		return strconv.Itoa(x)
	case float32:
		return strconv.FormatFloat(float64(x), 'f', -1, 32)
	// ...similar cases for int16, uint32, and so on...
	case bool:
		if x {
			return "true"
		}
		return "false"
	default:
		// array, chan, func, map, pointer, slice, struct
		return "???"
	}
}

func main() {
	var i int = 1
	// var s string = "hello"
	var f float32 = 1.11
	log.Println(Sprint(i))
	log.Println(Sprint(f))
}
