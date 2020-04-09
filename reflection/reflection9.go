package main

import (
	"fmt"
	"reflect"
)

type User struct {
	Name string
}

func main() {
	// T => *T
	u1 := User{"a8m"}
	p1 := ptr(reflect.ValueOf(u1))
	fmt.Println(u1 == p1.Elem().Interface())

	// *T => **T
	u2 := &User{"a8m"}
	p2 := ptr(reflect.ValueOf(u2))
	fmt.Println(*u2 == p2.Elem().Elem().Interface())
}

// ptr wraps the given value with pointer: V => *V, *V => **V, etc.
func ptr(v reflect.Value) reflect.Value {
	pt := reflect.PtrTo(v.Type()) // create a *T type.
	pv := reflect.New(pt.Elem())  // create a reflect.Value of type *T.
	pv.Elem().Set(v)              // sets pv to point to underlying value of v.
	return pv
}