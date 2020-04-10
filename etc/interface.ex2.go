package main

import (
	"fmt"
	"reflect"
)

// func isNil(i interface{}) bool {
// 	return i == nil || reflect.ValueOf(i).IsNil()
// }

func isNil(i interface{}) bool {
	if i == nil {
		return true
	}
	switch reflect.TypeOf(i).Kind() {
	case reflect.Ptr, reflect.Map, reflect.Array, reflect.Chan, reflect.Slice:
		return reflect.ValueOf(i).IsNil()
	}
	return false
}

type Animal interface {
	MakeSound() string
}
type Dog struct{}

func (d *Dog) MakeSound() string {
	return "Bark"
}

func main() {
	var d *Dog = nil
	var a Animal = d
	fmt.Println(isNil(a))
}

type Cat struct{}

func (c Cat) MakeSound() string {
	return "Meow"
}

// func main() {
// 	var c Cat
// 	var a Animal = c
// 	fmt.Println(isNil(a))
// }
