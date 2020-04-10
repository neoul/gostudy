package main

import (
	"fmt"
	"reflect"
)

type Config struct {
    Name string
    Meta struct {
        Desc string
        Properties map[string]string
        Users []string
    }
}

func initializeStruct(t reflect.Type, v reflect.Value) {
  for i := 0; i < v.NumField(); i++ {
    f := v.Field(i)
    ft := t.Field(i)
    switch ft.Type.Kind() {
    case reflect.Map:
      f.Set(reflect.MakeMap(ft.Type))
    case reflect.Slice:
      f.Set(reflect.MakeSlice(ft.Type, 0, 0))
    case reflect.Chan:
      f.Set(reflect.MakeChan(ft.Type, 0))
    case reflect.Struct:
      initializeStruct(ft.Type, f)
    case reflect.Ptr:
      fv := reflect.New(ft.Type.Elem())
      initializeStruct(ft.Type.Elem(), fv.Elem())
      f.Set(fv)
    default:
    }
  }
}

func main() {
	z := []int{1,2, 3}
	y := &z
	fmt.Println(append(z, 4))
	fmt.Println(z)
	fmt.Println(append(*y, 4))
	*y = append(*y, 4)
	fmt.Println(*y)
    // one way is to have a value of the type you want already
    a := 1
    // reflect.New works kind of like the built-in function new
    // We'll get a reflected pointer to a new int value
    intPtr := reflect.New(reflect.TypeOf(a))
    // Just to prove it
    b := intPtr.Elem().Interface().(int)
    // Prints 0
    fmt.Println(b)

    // We can also use reflect.New without having a value of the type
    var nilInt *int
	intType := reflect.TypeOf(nilInt).Elem()
    intPtr2 := reflect.New(intType)
    // Same as above
	x := intPtr2.Elem().Interface().(int)
	// fmt.Printf("%T", x)
    // Prints 0 again
	fmt.Println(x)
	
    t := reflect.TypeOf(Config{})
    v := reflect.New(t)
    initializeStruct(t, v.Elem())
    c := v.Interface().(*Config)
    c.Meta.Properties["color"] = "red" // map was already made!
    c.Meta.Users = append(c.Meta.Users, "srid") // so was the slice.
    fmt.Println(v.Interface())
}