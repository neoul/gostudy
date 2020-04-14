package model // import "github.com/neoul/gostudy/datablock/model"

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

type Lookup interface {
	Lookup(string) string
}

type DefaultValSet string

func (dv *DefaultValSet) Lookup(key string) string {
	values := string(*dv)
	if len(key) == 0 {
		return values
	}
	for _, kv := range strings.Split(values, ",") {
		s := strings.Split(kv, "=")
		if len(s) == 2 && s[0] == key {
			return strings.Trim(s[1], " ")
		}
	}
	return ""
}

// GetVarString - get the string of a variable value.
func GetVarString(src interface{}) string {
	v := reflect.ValueOf(src)
	t := reflect.TypeOf(src)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
		if v.IsNil() {
			// v = reflect.New(t)
			return "nil"
		}
		v = v.Elem()
	}
	switch t.Kind() {
	case reflect.String:
		return v.String()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return fmt.Sprintf("%d", v.Int())
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return fmt.Sprintf("%d", v.Uint())
	case reflect.Float32, reflect.Float64:
		return fmt.Sprintf("%f", v.Float())
	case reflect.Bool:
		if v.Bool() {
			return "true"
		}
		return "false"
	default:
		return ""
	}
}


// IsConcreteType reports whether t is a concrete type (built-in scalar type)
func IsConcreteType(t reflect.Type) bool {
	if t == reflect.TypeOf(nil) {
		return false
	}
	switch t.Kind() {
	case reflect.Array, reflect.Chan, reflect.Func, reflect.Interface, reflect.Map, reflect.Slice, reflect.Struct:
		return false
	case reflect.UnsafePointer, reflect.Complex64, reflect.Complex128: // not support
		return false
	case reflect.Ptr:
		return IsConcreteType(t.Elem())
	default:
		return true
	}
}



// initVar initializes the pointed variable to strval according to the type of the variable.
func initVar(t reflect.Type, v reflect.Value, strname string, strval string) error {
	if t.Kind() != reflect.Ptr {
		return fmt.Errorf("no ptr variable %s=%s(%s)", strname, t.Kind(), strval)
	}
	switch t.Elem().Kind() {
	case reflect.String:
		v.Elem().SetString(strval)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		val, err := strconv.ParseInt(strval, 10, 64)
		fmt.Printf("T=%T\n", val)
		fmt.Println(v.Elem().CanSet())
		if err != nil {
			return fmt.Errorf("Parse%s %s=%s", t.Elem().Kind(), strname, strval)
		}
		if v.Elem().OverflowInt(val) {
			return fmt.Errorf("Overflow %s=%s", strname, strval)
		}
		v.Elem().SetInt(val)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		val, err := strconv.ParseUint(strval, 10, 64)
		if err != nil {
			return fmt.Errorf("Parse%s %s=%s", t.Elem().Kind(), strname, strval)
		}
		if v.Elem().OverflowUint(val) {
			return fmt.Errorf("Overflow %s=%s", strname, strval)
		}
		v.Elem().SetUint(val)
		
	case reflect.Float32, reflect.Float64:
		val, err := strconv.ParseFloat(strval, 64)
		if err != nil {
			return fmt.Errorf("Parse%s %s=%s", t.Elem().Kind(), strname, strval)
		}
		if v.Elem().OverflowFloat(val) {
			return fmt.Errorf("Overflow %s=%s", strname, strval)
		}
		v.Elem().SetFloat(val)
	case reflect.Bool:
		var boolvar bool
		if strval == "true" || strval == "True" {
			boolvar = true
		} else if strval == "false" || strval == "False" {
			boolvar = false
		} else {
			return fmt.Errorf("Parse%s %s=%s", t.Elem().Kind(), strname, strval)
		}
		v.Elem().SetBool(boolvar)
	}
	return nil
}

func getChildName(ft reflect.StructField) (string, string) {
	prefix := ""
	name := ft.Name
	tag := ft.Tag.Get("json")
	
	if tag != "" {
		name = tag
	}
	tag = ft.Tag.Get("yaml")
	if tag != "" {
		name = tag
	}
	tag = ft.Tag.Get("path")
	if tag != "" {
		name = tag
	}
	prefix = ft.Tag.Get("module")
	// fmt.Printf("%s:%s\n", prefix, name)
	return prefix, name
}

func initStruct(t reflect.Type, v reflect.Value, values Lookup) error {
	fmt.Println("struct-type:", t)
	if !v.IsValid() || !v.CanSet() {
		return fmt.Errorf("Not assignable (%s)", v)
	}
	for i := 0; i < v.NumField(); i++ {
		
		fv := v.Field(i)
		ft := t.Field(i)
		if !fv.IsValid() || !fv.CanSet() {
			continue
		}
		fmt.Println(" - field-type:", ft)
		_, name := getChildName(ft)
		switch ft.Type.Kind() {
		case reflect.Map:
			fv.Set(reflect.MakeMap(ft.Type))
		case reflect.Slice:
			fv.Set(reflect.MakeSlice(ft.Type, 0, 0))
		case reflect.Chan:
			fv.Set(reflect.MakeChan(ft.Type, 0))
		case reflect.Struct:
			err := initStruct(ft.Type, fv, values)
			if err != nil {
				return err
			}
		case reflect.Ptr:
			if IsConcreteType(ft.Type) {
				strval := values.Lookup(name)
				rv := reflect.New(ft.Type.Elem())
				err := initVar(ft.Type, rv, name, strval)
				if err != nil {
					return err
				}
				fv.Set(rv)
			} else {
				newfv := reflect.New(ft.Type.Elem())
				err := initStruct(ft.Type.Elem(), newfv.Elem(), values)
				if err != nil {
					return err
				}
				fmt.Println(newfv)
				fv.Set(newfv)
			}
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			strval := values.Lookup(name)
			val, err := strconv.ParseInt(strval, 10, 64)
			if err != nil {
				return fmt.Errorf("Parse%s %s=%s", ft.Type.Kind(), name, strval)
			}
			if fv.OverflowInt(val) {
				return fmt.Errorf("Overflow %s=%s", name, strval)
			}
			fv.SetInt(val)
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
			strval := values.Lookup(name)
			val, err := strconv.ParseUint(strval, 10, 64)
			if err != nil {
				return fmt.Errorf("Parse%s %s=%s", ft.Type.Kind(), name, strval)
			}
			if fv.OverflowUint(val) {
				return fmt.Errorf("Overflow %s=%s", name, strval)
			}
			fv.SetUint(val)
		case reflect.Float32, reflect.Float64:
			strval := values.Lookup(name)
			val, err := strconv.ParseFloat(strval, 64)
			if err != nil {
				return fmt.Errorf("Parse%s %s=%s", ft.Type.Kind(), name, strval)
			}
			if fv.OverflowFloat(val) {
				return fmt.Errorf("Overflow %s=%s", name, strval)
			}
			fv.SetFloat(val)
		case reflect.Bool:
			strval := values.Lookup(name)
			var boolvar bool = false
			if strval == "true" || strval == "True" {
				boolvar = true
			} else if strval == "false" || strval == "False" {
				boolvar = false
			} else {
				return fmt.Errorf("Parse%s %s=%s", ft.Type.Kind(), name, strval)
			}
			fv.SetBool(boolvar)
		case reflect.String:
			strval := values.Lookup(name)
			fv.SetString(strval)
		default:
			return fmt.Errorf("Not supported type (%s) in %s", ft.Type.Kind(), name)
		}
		fmt.Println(" - field-value:", fv)
	}
	// fmt.Println(v)
	fmt.Println("struct-value:", v)
	return nil
}

// Set - Set all structure fields of the top.
func Set(top interface{}, value Lookup) (interface{}, error) {
	v := reflect.ValueOf(top)
	t := reflect.TypeOf(top)
	if v.Kind() != reflect.Ptr {
		return top, fmt.Errorf("no ptr variable (%s)", reflect.TypeOf(top))
	}
	if IsConcreteType(t) {
		if v.IsNil() {
			v = reflect.New(t.Elem())
		}
		initVar(t, v, t.Name(), value.Lookup(""))
	} else {
		// fmt.Println("nil?", v.IsNil(), t, v)
		t = t.Elem()
		if v.IsNil() {
			v = reflect.New(t)
		}
		initStruct(t, v.Elem(), value)
		// fmt.Println("struct-value:", v)
	}
	return v.Interface(), nil
}


func PrintStruct(t reflect.Type, v reflect.Value, indent string) {
	fmt.Printf(indent)
	fmt.Printf("%s(%s):\n", t.Name(), t.Kind())
	indent += " . "
	for i := 0; i < v.NumField(); i++ {
		fv := v.Field(i)
		ft := t.Field(i)
		if !fv.IsValid() || !fv.CanSet() {
			continue
		}
		_, name := getChildName(ft)
		fmt.Printf(indent)
		switch ft.Type.Kind() {
		case reflect.Map:
			fmt.Printf("%s(%s):", name, ft.Type.Kind())
		case reflect.Slice:
			fmt.Printf("%s(%s):", name, ft.Type.Kind())
		case reflect.Chan:
			fmt.Printf("%s(%s):", name, ft.Type.Kind())
		case reflect.Struct:
			PrintStruct(ft.Type, fv, indent)
		case reflect.Ptr:
			if IsConcreteType(ft.Type) {
				fmt.Printf("%s: %s\n", name, GetVarString(fv.Interface()))
			} else {
				PrintStruct(ft.Type, fv, indent)
			}
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			fmt.Printf("%s(%s): %d", name, ft.Type.Kind(), fv.Int())
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
			fmt.Printf("%s(%s): %d", name, ft.Type.Kind(), fv.Uint())
		case reflect.Float32, reflect.Float64:
			fmt.Printf("%s(%s): %f", name, ft.Type.Kind(), fv.Float())
		case reflect.Bool:
			fmt.Printf("%s(%s): %t", name, ft.Type.Kind(), fv.Bool())
		case reflect.String:
			fmt.Printf("%s(%s): %s", name, ft.Type.Kind(), fv.String())
		default:
			fmt.Printf("%s(%s): %s", name, ft.Type.Kind(), "?")
		}
		fmt.Printf("\n")
	}
}

func Print(top interface{}) {
	var indent string
	v := reflect.ValueOf(top)
	t := reflect.TypeOf(top)
	if v.Kind() != reflect.Ptr {
		return
	}
	if IsConcreteType(t) {
		fmt.Printf("%s(%s): %s\n", t.Name(), t.Kind(), GetVarString(top))
	} else {
		t = t.Elem()
		if v.IsNil() {
			fmt.Printf("%s(%s): %s\n", t.Name(), t.Kind(), GetVarString(top))
		} else {
			PrintStruct(t, v.Elem(), indent)
		}
	}
}