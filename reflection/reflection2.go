package main
// https://github.com/a8m/reflect-examples
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
	u := &User{Name: "Ariel Mashraki"}
	// Elem returns the value that the pointer u points to.
	v := reflect.ValueOf(u).Elem()
	f := v.FieldByName("Github")
	// make sure that this field is defined, and can be changed.
	if !f.IsValid() || !f.CanSet() {
		return
	}
	if f.Kind() != reflect.String || f.String() != "" {
		return
	}
	f.SetString("a8m")
	fmt.Printf("Github username was changed to: %q\n", u.Github)
}