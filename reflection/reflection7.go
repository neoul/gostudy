package main

import (
	"fmt"
	"reflect"
	"strings"
)

type User struct {
	Email   string `kv:"email,omitempty"`
	Name    string `kv:"name,omitempty"`
	Github  string `kv:"github,omitempty"`
	private string
}

func main() {
	var (
		u = User{Name: "Ariel", Github: "a8m"}
		v = struct {
			A, B, C string
		}{
			"foo",
			"bar",
			"baz",
		}
		w = &User{}
	)
	fmt.Println(encode(u))
	fmt.Println(encode(v))
	fmt.Println(encode(w))
}

// this example supports only structs, and assume their
// fields are type string.
func encode(i interface{}) (string, error) {
	v := reflect.ValueOf(i)
	t := v.Type()
	if t.Kind() != reflect.Struct {
		return "", fmt.Errorf("type %s is not supported", t.Kind())
	}
	var s []string
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		// skip unexported fields. from godoc:
		// PkgPath is the package path that qualifies a lower case (unexported)
		// field name. It is empty for upper case (exported) field names.
		if f.PkgPath != "" {
			continue
		}
		fv := v.Field(i)
		key, omit := readTag(f)
		// skip empty values when "omitempty" set.
		if omit && fv.String() == "" {
			continue
		}
		s = append(s, fmt.Sprintf("%s=%s", key, fv.String()))
	}
	return strings.Join(s, ","), nil
}

func readTag(f reflect.StructField) (string, bool) {
	val, ok := f.Tag.Lookup("kv")
	if !ok {
		return f.Name, false
	}
	opts := strings.Split(val, ",")
	omit := false
	if len(opts) == 2 {
		omit = opts[1] == "omitempty"
	}
	return opts[0], omit
}