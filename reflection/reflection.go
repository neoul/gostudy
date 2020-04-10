package main

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
				strval := values.Lookup(ft.Name)
				rv := reflect.New(ft.Type.Elem())
				err := initVar(ft.Type, rv, ft.Name, strval)
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
			strval := values.Lookup(ft.Name)
			val, err := strconv.ParseInt(strval, 10, 64)
			if err != nil {
				return fmt.Errorf("Parse%s %s=%s", ft.Type.Kind(), ft.Name, strval)
			}
			if fv.OverflowInt(val) {
				return fmt.Errorf("Overflow %s=%s", ft.Name, strval)
			}
			fv.SetInt(val)
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
			strval := values.Lookup(ft.Name)
			val, err := strconv.ParseUint(strval, 10, 64)
			if err != nil {
				return fmt.Errorf("Parse%s %s=%s", ft.Type.Kind(), ft.Name, strval)
			}
			if fv.OverflowUint(val) {
				return fmt.Errorf("Overflow %s=%s", ft.Name, strval)
			}
			fv.SetUint(val)
		case reflect.Float32, reflect.Float64:
			strval := values.Lookup(ft.Name)
			val, err := strconv.ParseFloat(strval, 64)
			if err != nil {
				return fmt.Errorf("Parse%s %s=%s", ft.Type.Kind(), ft.Name, strval)
			}
			if fv.OverflowFloat(val) {
				return fmt.Errorf("Overflow %s=%s", ft.Name, strval)
			}
			fv.SetFloat(val)
		case reflect.Bool:
			strval := values.Lookup(ft.Name)
			var boolvar bool = false
			if strval == "true" || strval == "True" {
				boolvar = true
			} else if strval == "false" || strval == "False" {
				boolvar = false
			} else {
				return fmt.Errorf("Parse%s %s=%s", ft.Type.Kind(), ft.Name, strval)
			}
			fv.SetBool(boolvar)
		case reflect.String:
			strval := values.Lookup(ft.Name)
			fv.SetString(strval)
		default:
			return fmt.Errorf("Not supported type (%s) in %s", ft.Type.Kind(), ft.Name)
		}
		fmt.Println(" - field-value:", fv)
	}
	// fmt.Println(v)
	fmt.Println("struct-value:", v)
	return nil
}

// Set and initialize all types of variables with value or zero.
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


type Site struct {
	Address *string
	Hit int
}

type Friend struct {
	FName string
}

type User struct {
	Name    string
	Github  string
	PString  *string
	private string
	Site Site
	SitePtr *Site
	Friend map[string]*Friend
	FriendPtr map[string]Friend
	Item []*string
}

func main() {
	var (
		v0 User
		v1 *User
		v2 = new(User)
		v3 struct{ Name string }
		// dv = DefaultValSet()
		d = "Name=Ariel,Github=a8m,PString=pointer-value,Hit=10,Address=github.com/neoul"
	)
	dv := DefaultValSet(d)
	Set(&v0, &dv)
	v0.Friend["Boy"]=&Friend{FName: "friend1"}
	fmt.Println("+++", v0, v0.Friend["Boy"])
	x, _ := Set(v1, &dv)
	v1 = x.(*User)
	fmt.Println("+++", v1)
	Set(v2, &dv)
	fmt.Println("+++", v2)
	Set(&v3, &dv)
	fmt.Println("+++", v3)
	val1 := 1
	dv = DefaultValSet("10")
	Set(&val1, &dv)
	fmt.Println(val1)
}

