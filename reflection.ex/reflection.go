package main

import (
	"fmt"
	"reflect"
)

type Person struct {
	name string `tag1:"1" tag2:"2"`
	age  int    `tag1:"나이" tag2:"Age"`
}

func main() {
	var i int = 1
	var s string = "hello world"
	var f float32 = 1.3

	fmt.Println(reflect.TypeOf(i))
	fmt.Println(reflect.TypeOf(s))
	fmt.Println(reflect.TypeOf(f))
	t := reflect.TypeOf(f)
	v := reflect.ValueOf(f)
	fmt.Println("")
	fmt.Println("float32 reflection")
	fmt.Println("==================")
	fmt.Println(t.Name())
	fmt.Println(t.Size())
	fmt.Println(t.Kind() == reflect.Float32)
	fmt.Println(t.Kind() == reflect.Float64)
	fmt.Println(v.Type())
	fmt.Println(v.Kind() == reflect.Float32)
	fmt.Println(v.Kind() == reflect.Float64)
	fmt.Println(v.Float())
	fmt.Println(v)

	fmt.Println("")
	fmt.Println("struct reflection")
	fmt.Println("==================")
	var d Person = Person{"myname", 3}
	var p *Person = &d
	fmt.Println(reflect.TypeOf(d))
	name, ok := reflect.TypeOf(d).FieldByName("name")
	fmt.Println("num of fields", reflect.TypeOf(d).NumField())
	fmt.Println(ok, name.Tag.Get("tag1"), name.Tag.Get("tag2"))
	age, ok := reflect.TypeOf(d).FieldByName("age")
	fmt.Println(ok, age.Tag.Get("tag1"), age.Tag.Get("tag2"))
	fmt.Println(reflect.TypeOf(p))
	fmt.Println(reflect.ValueOf(p))
	fmt.Println(reflect.ValueOf(p).Elem()) // reflection of pointer
	fmt.Println(reflect.ValueOf(p).Elem().FieldByName("name"),
		reflect.ValueOf(p).Elem().FieldByName("age"))

	fmt.Println("")
	fmt.Println("interface reflection")
	fmt.Println("==================")
	var b interface{}
	b = 1
	fmt.Println(reflect.TypeOf(b))
	fmt.Println(reflect.ValueOf(b))
	fmt.Println(reflect.ValueOf(b).Int())
	// fmt.Println(reflect.ValueOf(b).Elem()) // Runtime error
}
