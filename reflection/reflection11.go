package main

import (
	"fmt"
	"reflect"
)

type Say interface {
	Say()
}

type xx struct {
	Say  Say
	Name string
}

type yy struct {
	Name string
}

func (y *yy) Say() {
	fmt.Println("Say")
}

func main() {
	// Data of interface type
	s := xx{Name: "name"}
	v := reflect.ValueOf(s)
	fmt.Println(v.Kind())

	fmt.Println(v.FieldByName("Name").Kind())
	fmt.Println(v.FieldByName("Say").Kind())
	fmt.Println(v.FieldByName("Say").Type())
	t := v.FieldByName("Say").Type()
	fmt.Println(t)
	fmt.Println(t.Implements)
	fmt.Println(v.FieldByName("Say").Elem())
	fmt.Println(v.FieldByName("Say").IsValid())
	fmt.Println(v.FieldByName("Say").IsZero())
	fmt.Println(v.FieldByName("Say").IsNil())

	s = xx{Name: "name", Say: &(yy{})}
	v = reflect.ValueOf(s)
	fmt.Println(v.Kind())

	fmt.Println(v.FieldByName("Name").Kind())
	fmt.Println(v.FieldByName("Say").Kind())
	fmt.Println(v.FieldByName("Say").Type())
	fmt.Println(v.FieldByName("Say").Elem())
	fmt.Println(v.FieldByName("Say").IsValid())
	fmt.Println(v.FieldByName("Say").IsZero())
	fmt.Println(v.FieldByName("Say").IsNil())
	fmt.Println(v.FieldByName("Say").Elem().Type())
}
