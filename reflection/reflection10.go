package main

import (
	"fmt"
	"reflect"
	"strconv"
)

func convertConcreteType(src interface{}, dest interface{}) (interface{}, bool) {
	fmt.Printf("src(%T) -> dest(%T): %v\n", src, dest, src)
	st := reflect.TypeOf(src)
	sv := reflect.ValueOf(src)
	dt := reflect.TypeOf(dest)
	if dt.Kind() == reflect.String {
		switch st.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			return fmt.Sprintf("%d", sv.Int()), true
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
			return fmt.Sprintf("%d", sv.Uint()), true
		case reflect.Float32, reflect.Float64:
			return fmt.Sprintf("%f", sv.Float()), true
		case reflect.Bool:
			if sv.Bool() {
				return "true", true
			}
			return "false", true
		}
	}
	if st.Kind() == reflect.String {
		rv := reflect.New(dt).Elem()
		dv := reflect.ValueOf(dest)
		switch dt.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			val, err := strconv.ParseInt(src.(string), 10, 64)
			if err != nil {
				return dest, false
			}
			if dv.OverflowInt(val) {
				return dest, false
			}
			rv.SetInt(val)
			return rv.Interface(), true
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
			val, err := strconv.ParseUint(src.(string), 10, 64)
			if err != nil {
				return dest, false
			}
			if dv.OverflowUint(val) {
				return dest, false
			}
			rv.SetUint(val)
			return rv.Interface(), true
		case reflect.Float32, reflect.Float64:
			val, err := strconv.ParseFloat(src.(string), 64)
			if err != nil {
				return dest, false
			}
			if dv.OverflowFloat(val) {
				return dest, false
			}
			rv.SetFloat(val)
			return rv.Interface(), true
		case reflect.Bool:
			s := src.(string)
			if s == "true" || s == "True" || s == "TRUE" {
				rv.SetBool(true)
			} else {
				rv.SetBool(false)
			}
			return rv.Interface(), true
		}
	}
	if dt.Kind() == reflect.Bool {
		switch st.Kind() {
		case reflect.String:
			if len(sv.String()) > 0 {
				return true, true
			}
			return false, true
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			if sv.Int() != 0 {
				return true, true
			}
			return false, true
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
			if sv.Uint() != 0 {
				return true, true
			}
			return false, true
		case reflect.Float32, reflect.Float64:
			if sv.Float() != 0 {
				return true, true
			}
			return false, true
		}
	}
	if st.Kind() == reflect.Bool {
		rv := reflect.New(dt).Elem()
		switch dt.Kind() {
		case reflect.String:
			if sv.Bool() {
				return "true", true
			}
			return "false", true
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			if sv.Bool() {
				rv.SetInt(1)
			} else {
				rv.SetInt(0)
			}
			return rv.Interface(), true
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
			if sv.Bool() {
				rv.SetUint(1)
			} else {
				rv.SetUint(0)
			}
			return rv.Interface(), true
		case reflect.Float32, reflect.Float64:
			if sv.Bool() {
				rv.SetFloat(1)
			} else {
				rv.SetFloat(0)
			}
			return rv.Interface(), true
		}
	}

	if st.ConvertibleTo(dt) {
		val := sv.Convert(dt)
		return val.Interface(), true
	}
	rv := reflect.New(dt).Elem()
	return rv.Interface(), false
}



func main() {
	var a int32 = 10
	var b bool = true
	var c uint64 = 12
	var d string = "10"
	var e float32 = 1
	fmt.Println(convertConcreteType(a, b))
	fmt.Println(convertConcreteType(b, c))
	fmt.Println(convertConcreteType(c, d))
	fmt.Println(convertConcreteType(e, a))
	fmt.Println(convertConcreteType(a, c))
	fmt.Println(convertConcreteType(b, d))
	fmt.Println(convertConcreteType(c, e))
	fmt.Println(convertConcreteType(d, a))
	fmt.Println(convertConcreteType(e, a))
	fmt.Println(convertConcreteType(a, d))
	fmt.Println(convertConcreteType(b, e))
	fmt.Println(convertConcreteType(c, a))
	fmt.Println(convertConcreteType(d, b))
	fmt.Println(convertConcreteType(e, c))
	fmt.Println(convertConcreteType(a, e))
	fmt.Println(convertConcreteType(b, a))
	fmt.Println(convertConcreteType(c, b))
	fmt.Println(convertConcreteType(d, c))
	fmt.Println(convertConcreteType(e, d))
	fmt.Println(convertConcreteType(d, d))
}
