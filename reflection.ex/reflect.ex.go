package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"reflect"
)

func main() {
	var w io.Writer = os.Stdout
	log.Println(reflect.TypeOf(w))
	log.Printf("%T", w) // Equal to TypeOf
	log.Println(reflect.ValueOf(w))
	log.Printf("%v", w) // Equal to ValueOf

	var emtpy interface{}
	emtpy = 10
	log.Println(emtpy.(int))
	log.Println(reflect.ValueOf(emtpy).Int())
	k := reflect.ValueOf(emtpy).Kind()
	
	log.Println(k)

	v := reflect.ValueOf(3) // a reflect.Value
	x := v.Interface()      // an interface{}
	i := x.(int)            // an int
	fmt.Printf("%d\n", i)   // "3"
	fmt.Println(x)
}
