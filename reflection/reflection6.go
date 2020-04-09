package main

import (
	"fmt"
	"reflect"
	"strings"
)

type User struct {
	Name    string
	Github  string
	private string
}

func main() {
	var (
		v0 User
		v1 *User
		v2 = new(User)
		v3 struct{ Name string }
		s  = "Name=Ariel,Github=a8m"
	)
	fmt.Println(decode(s, &v0), v0) // pass
	fmt.Println(decode(s, v1), v1)  // fail
	fmt.Println(decode(s, v2), v2)  // pass
	fmt.Println(decode(s, v3), v3)  // fail
	fmt.Println(decode(s, &v3), v3) // pass
}

func decode(s string, i interface{}) error {
	v := reflect.ValueOf(i)
	if v.Kind() != reflect.Ptr || v.IsNil() {
		return fmt.Errorf("decode requires non-nil pointer")
	}
	// get the value that the pointer v points to.
	v = v.Elem()
	// assume that the input is valid.
	for _, kv := range strings.Split(s, ",") {
		s := strings.Split(kv, "=")
		f := v.FieldByName(s[0])
		// make sure that this field is defined, and can be changed.
		if !f.IsValid() || !f.CanSet() {
			continue
		}
		// assume all the fields are type string.
		f.SetString(s[1])
	}
	return nil
}