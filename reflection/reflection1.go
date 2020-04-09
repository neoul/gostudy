package main

import (
	"fmt"
	"reflect"
)

type User struct {
	Email  string `mcl:"email"`
	Name   string `mcl:"name"`
	Age    int    `mcl:"age"`
	Github string `mcl:"github" default:"a8m"`
}

func main() {
	var u interface{} = User{}
	// TypeOf returns the reflection Type that represents the dynamic type of u.
	t := reflect.TypeOf(u)
	// Kind returns the specific kind of this type. 
	if t.Kind() != reflect.Struct {
		return
	}
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		fmt.Println("mcl:", f.Tag.Get("mcl"), "default:", f.Tag.Get("default"))
	}
}