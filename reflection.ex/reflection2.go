package main

import (
	"fmt"
	"log"
	"reflect"
)

type Foo struct {
	FirstName string `tag_name:"tag 1"`
	LastName  string `tag_name:"tag 2"`
	Age       int    `tag_name:"tag 3"`
}

type map_str_interface map[string]interface{}

func (f *Foo) reflect() {
	val := reflect.ValueOf(f).Elem()

	for i := 0; i < val.NumField(); i++ {
		valueField := val.Field(i)
		typeField := val.Type().Field(i)
		tag := typeField.Tag

		log.Printf("Field Name: %s,\t Field Value: %v,\t Tag Value: %s\n", typeField.Name, valueField.Interface(), tag.Get("tag_name"))
	}
}

func inquiry_struct(s interface{}) {
	t := reflect.TypeOf(s)
	// v := reflect.ValueOf(s)
	// k := reflect.Kind(s)
	log.Println(t)
	// log.Println(reflect.Kind(v))

}

func isNil(i interface{}) bool {
	return i == nil || reflect.ValueOf(i).IsNil()
}

func main() {
	var i int = 1
	var s string = "hello"
	var f float32 = 1.0
	var m map_str_interface = make(map_str_interface)

	inquiry_struct(i)
	inquiry_struct(s)
	inquiry_struct(f)
	inquiry_struct(m)

	var v interface{}
	var t = reflect.ValueOf(&v).Type().Elem()
	fmt.Println(t.Kind() == reflect.Interface)

	// f := &Foo{
	// 	FirstName: "Drew",
	// 	LastName:  "Olson",
	// 	Age:       30,
	// }

	// f.reflect()
}
