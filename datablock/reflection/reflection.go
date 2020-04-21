package reflection // import "github.com/neoul/gostudy/datablock/reflection"

import (
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"

	"github.com/neoul/gostudy/datablock/log"
)

// Generate rule to create the example structs:
//go:generate go run ../../../github.com/openconfig/ygot/generator/generator.go -path=yang -output_file=object/example.go -package_name=object -generate_fakeroot -fakeroot_name=device yang/example.yang

var (
	ylog log.Log
)

func init() {
	ylog = log.NewLog("ydb2go", os.Stdout)
}

// GetPureType returns not reflect.Ptr type.
func GetPureType(t reflect.Type) reflect.Type {
	for ; t.Kind() == reflect.Ptr; t = t.Elem() {
	}
	return t
}

// IsTypeDeep reports whether t is k type.
func IsTypeDeep(t reflect.Type, kinds ...reflect.Kind) bool {
	for ; t.Kind() == reflect.Ptr; t = t.Elem() {
	}
	for _, k := range kinds {
		if t.Kind() == k {
			return true
		}
	}
	return false
}

// IsReferenceType returns true if t is a map, slice or channel
func IsReferenceType(t reflect.Type) bool {
	switch t.Kind() {
	case reflect.Slice, reflect.Chan, reflect.Map:
		return true
	}
	return false
}

// IsTypeStruct reports whether t is a struct type.
func IsTypeStruct(t reflect.Type) bool {
	return t.Kind() == reflect.Struct
}

// IsTypeStructPtr reports whether v is a struct ptr type.
func IsTypeStructPtr(t reflect.Type) bool {
	if t == reflect.TypeOf(nil) {
		return false
	}
	return t.Kind() == reflect.Ptr && t.Elem().Kind() == reflect.Struct
}

// IsTypeSlice reports whether v is a slice type.
func IsTypeSlice(t reflect.Type) bool {
	return t.Kind() == reflect.Slice
}

// IsTypeSlicePtr reports whether v is a slice ptr type.
func IsTypeSlicePtr(t reflect.Type) bool {
	if t == reflect.TypeOf(nil) {
		return false
	}
	return t.Kind() == reflect.Ptr && t.Elem().Kind() == reflect.Slice
}

// IsTypeMap reports whether v is a map type.
func IsTypeMap(t reflect.Type) bool {
	if t == reflect.TypeOf(nil) {
		return false
	}
	return t.Kind() == reflect.Map
}

// IsTypeInterface reports whether v is an interface.
func IsTypeInterface(t reflect.Type) bool {
	if t == reflect.TypeOf(nil) {
		return false
	}
	return t.Kind() == reflect.Interface
}

// IsTypeSliceOfInterface reports whether v is a slice of interface.
func IsTypeSliceOfInterface(t reflect.Type) bool {
	if t == reflect.TypeOf(nil) {
		return false
	}
	return t.Kind() == reflect.Slice && t.Elem().Kind() == reflect.Interface
}

// IsTypeChan reports whether v is a slice type.
func IsTypeChan(t reflect.Type) bool {
	return t.Kind() == reflect.Chan
}

// IsNilOrInvalidValue reports whether v is nil or reflect.Zero.
func IsNilOrInvalidValue(v reflect.Value) bool {
	return !v.IsValid() || (v.Kind() == reflect.Ptr && v.IsNil()) || IsValueNil(v.Interface())
}

// IsValueNil returns true if either value is nil, or has dynamic type {ptr,
// map, slice} with value nil.
func IsValueNil(value interface{}) bool {
	if value == nil {
		return true
	}
	switch reflect.TypeOf(value).Kind() {
	case reflect.Slice, reflect.Ptr, reflect.Map:
		return reflect.ValueOf(value).IsNil()
	}
	return false
}

// IsValueNilOrDefault returns true if either IsValueNil(value) or the default
// value for the type.
func IsValueNilOrDefault(value interface{}) bool {
	if IsValueNil(value) {
		return true
	}
	if !IsValueScalar(reflect.ValueOf(value)) {
		// Default value is nil for non-scalar types.
		return false
	}
	return value == reflect.New(reflect.TypeOf(value)).Elem().Interface()
}

// IsValuePtr reports whether v is a ptr type.
func IsValuePtr(v reflect.Value) bool {
	return v.Kind() == reflect.Ptr
}

// IsValueInterface reports whether v is an interface type.
func IsValueInterface(v reflect.Value) bool {
	return v.Kind() == reflect.Interface
}

// IsValueStruct reports whether v is a struct type.
func IsValueStruct(v reflect.Value) bool {
	return v.Kind() == reflect.Struct
}

// IsValueStructPtr reports whether v is a struct ptr type.
func IsValueStructPtr(v reflect.Value) bool {
	return v.Kind() == reflect.Ptr && IsValueStruct(v.Elem())
}

// IsValueMap reports whether v is a map type.
func IsValueMap(v reflect.Value) bool {
	return v.Kind() == reflect.Map
}

// IsValueSlice reports whether v is a slice type.
func IsValueSlice(v reflect.Value) bool {
	return v.Kind() == reflect.Slice
}

// IsValueScalar reports whether v is a scalar type.
func IsValueScalar(v reflect.Value) bool {
	if IsNilOrInvalidValue(v) {
		return false
	}
	if IsValuePtr(v) {
		if v.IsNil() {
			return false
		}
		v = v.Elem()
	}
	return !IsValueStruct(v) && !IsValueMap(v) && !IsValueSlice(v)
}

// ValuesAreSameType returns true if v1 and v2 has the same reflect.Type,
// otherwise it returns false.
func ValuesAreSameType(v1 reflect.Value, v2 reflect.Value) bool {
	return v1.Type() == v2.Type()
}

// IsValueInterfaceToStructPtr reports whether v is an interface that contains a
// pointer to a struct.
func IsValueInterfaceToStructPtr(v reflect.Value) bool {
	return IsValueInterface(v) && IsValueStructPtr(v.Elem())
}

// IsStructValueWithNFields returns true if the reflect.Value representing a
// struct v has n fields.
func IsStructValueWithNFields(v reflect.Value, n int) bool {
	return IsValueStruct(v) && v.NumField() == n
}

// IsSimpleType - true if built-in simple variable type
func IsSimpleType(t reflect.Type) bool {
	switch t.Kind() {
	case reflect.Ptr:
		return IsSimpleType(t.Elem())
	case reflect.Array, reflect.Chan, reflect.Func, reflect.Interface,
		reflect.Map, reflect.Slice, reflect.Struct,
		reflect.UnsafePointer, reflect.Complex64, reflect.Complex128:
		return false
	default:
		return true
	}
}

var maxValueStrLen = 150

// ValueStrDebug returns "<not calculated>" if the package global variable
// debugLibrary is not set. Otherwise, it is the same as ValueStr.
// Use this function instead of ValueStr for debugging purpose, e.g. when the
// output is passed to DbgPrint, because ValueStr calls can be the bottleneck
// for large input.
func ValueStrDebug(value interface{}) string {
	return ValueStr(value)
}

// ValueStr returns a string representation of value which may be a value, ptr,
// or struct type.
func ValueStr(value interface{}) string {
	out := valueStrInternal(value)
	if len(out) > maxValueStrLen {
		out = out[:maxValueStrLen] + "..."
	}
	return out
}

// ValueStrInternal is the internal implementation of ValueStr.
func valueStrInternal(value interface{}) string {
	v := reflect.ValueOf(value)
	kind := v.Kind()
	switch kind {
	case reflect.Ptr:
		if v.IsNil() || !v.IsValid() {
			return "nil"
		}
		return strings.Replace(ValueStr(v.Elem().Interface()), ")", " ptr)", -1)
	case reflect.Slice:
		var out string
		for i := 0; i < v.Len(); i++ {
			if i != 0 {
				out += ", "
			}
			out += ValueStr(v.Index(i).Interface())
		}
		return "[ " + out + " ]"
	case reflect.Struct:
		var out string
		for i := 0; i < v.NumField(); i++ {
			if i != 0 {
				out += ", "
			}
			if !v.Field(i).CanInterface() {
				continue
			}
			out += ValueStr(v.Field(i).Interface())
		}
		return "{ " + out + " }"
	}
	out := fmt.Sprintf("%v (%v)", value, kind)
	if len(out) > maxValueStrLen {
		out = out[:maxValueStrLen] + "..."
	}
	return out
}

func printValue(v reflect.Value, ptrcnt int) string {
	if IsNilOrInvalidValue(v) {
		return "(???)"
	}
	if IsValuePtr(v) {
		ptrcnt++
		return "*" + printValue(v.Elem(), ptrcnt)
	}
	s := fmt.Sprintf("%s(", v.Type())
	for i := 0; i < ptrcnt; i++ {
		s = s + "&"
	}
	s = s + fmt.Sprintf("%v)", v)
	return s
}

// PrintValue writes the value type and value to a string
func PrintValue(value interface{}) string {
	v := reflect.ValueOf(value)
	return printValue(v, 0)
}


func setMapValue(mv reflect.Value, key interface{}, value interface{}) error {
	if mv.Kind() == reflect.Ptr {
		if mv.IsNil() {
			cv := newValueMap(mv.Type())
			mv.Set(cv)
		} else {
			return setMapValue(mv.Elem(), key, value)
		}
	}
	mt := mv.Type()
	kt := mt.Key()
	kv := newValue(kt, key)
	vv := newValue(mt.Elem(), value)
	mv.SetMapIndex(kv, vv)
	return nil
}

func setSliceValue(sv reflect.Value, element interface{}) reflect.Value {
	st := sv.Type()
	if st.Kind() == reflect.Ptr {
		if st.Elem().Kind() == reflect.Slice {
			et := st.Elem().Elem()
			ev := newValue(et, element)
			ylog.Debug("::", ev)
			sv.Elem().Set(reflect.Append(sv.Elem(), ev))
			return sv
		} else if st.Elem().Kind() == reflect.Ptr {
			return setSliceValue(sv.Elem(), element)
		} else {
		}
		// if sv.IsNil() {
		// 	cv := newValueSlice(sv.Type())
		// 	sv.Set(cv)
		// } else {
		// 	return setSliceValue(sv.Elem(), element)
		// }
	}
	ev := newValue(st.Elem(), element)
	num := sv.Len()
	nslice := reflect.MakeSlice(st, num+1, num+1)
	reflect.Copy(nslice, sv)
	nslice.Index(num).Set(ev)
	return nslice
}

func setValueScalar(v reflect.Value, value interface{}) error {
	dv := v
	if dv.Kind() == reflect.Ptr {
		dv = v.Elem()
		if dv.Kind() == reflect.Ptr {
			if dv.IsNil() { // e.g. **type
				dv = reflect.New(dv.Type().Elem())
				ylog.Debug(PrintValue(dv.Interface()))
				v.Elem().Set(dv)
			}
			return setValueScalar(dv, value)
		}
	}
	dt := dv.Type()
	st := reflect.TypeOf(value)
	sv := reflect.ValueOf(value)
	if dt.Kind() == reflect.String {
		switch st.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			dv.SetString(fmt.Sprintf("%d", sv.Int()))
			return nil
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
			dv.SetString(fmt.Sprintf("%d", sv.Uint()))
			return nil
		case reflect.Float32, reflect.Float64:
			dv.SetString(fmt.Sprintf("%f", sv.Float()))
			return nil
		case reflect.Bool:
			dv.SetString(fmt.Sprint(sv.Bool()))
			return nil
		}
	}
	if st.Kind() == reflect.String {
		srcstring := value.(string)
		switch dt.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			if len(srcstring) == 0 {
				dv.SetInt(0)
			} else {
				val, err := strconv.ParseInt(srcstring, 10, 64)
				if err != nil {
					return err
				}
				if dv.OverflowInt(val) {
					return fmt.Errorf("overflowInt: %s", PrintValue(val))
				}
				dv.SetInt(val)
			}
			return nil
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
			if len(srcstring) == 0 {
				dv.SetUint(0)
			} else {
				val, err := strconv.ParseUint(srcstring, 10, 64)
				if err != nil {
					return err
				}
				if dv.OverflowUint(val) {
					return fmt.Errorf("OverflowUint: %s", PrintValue(val))
				}
				dv.SetUint(val)
			}
			return nil
		case reflect.Float32, reflect.Float64:
			if len(srcstring) == 0 {
				dv.SetFloat(0)
			} else {
				val, err := strconv.ParseFloat(srcstring, 64)
				if err != nil {
					return err
				}
				if dv.OverflowFloat(val) {
					return fmt.Errorf("OverflowFloat: %s", PrintValue(val))
				}
				dv.SetFloat(val)
			}
			return nil
		case reflect.Bool:
			if srcstring == "true" || srcstring == "True" || srcstring == "TRUE" {
				dv.SetBool(true)
			} else {
				dv.SetBool(false)
			}
			return nil
		}
	}
	if dt.Kind() == reflect.Bool {
		switch st.Kind() {
		case reflect.String:
			if len(sv.String()) > 0 && sv.String() == "true" || sv.String() == "True" || sv.String() == "TRUE" {
				dv.SetBool(true)
			} else {
				dv.SetBool(false)
			}
			return nil
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			if sv.Int() != 0 {
				dv.SetBool(true)
			} else {
				dv.SetBool(false)
			}
			return nil
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
			if sv.Uint() != 0 {
				dv.SetBool(true)
			} else {
				dv.SetBool(false)
			}
			return nil
		case reflect.Float32, reflect.Float64:
			if sv.Float() != 0 {
				dv.SetBool(true)
			} else {
				dv.SetBool(false)
			}
			return nil
		}
	}
	if st.Kind() == reflect.Bool {
		switch dt.Kind() {
		case reflect.String:
			if sv.Bool() {
				dv.SetString("true")
			} else {
				dv.SetString("false")
			}
			return nil
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			if sv.Bool() {
				dv.SetInt(1)
			} else {
				dv.SetInt(0)
			}
			return nil
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
			if sv.Bool() {
				dv.SetUint(1)
			} else {
				dv.SetUint(0)
			}
			return nil
		case reflect.Float32, reflect.Float64:
			if sv.Bool() {
				dv.SetFloat(1)
			} else {
				dv.SetFloat(0)
			}
			return nil
		}
	}
	if st.ConvertibleTo(dt) {
		dv.Set(sv.Convert(dt))
		return nil
	}
	return fmt.Errorf("Not Convertible: %s", PrintValue(v.Interface()))
}

func newValueStruct(t reflect.Type) reflect.Value {
	pv := reflect.New(t)
	pt := reflect.PtrTo(t)
	ylog.Debug(pv, pt)
	if pv.Elem().Kind() == reflect.Ptr {
		cv := newValueStruct(t.Elem())
		pv.Elem().Set(cv)
		return pv
	}
	pve := pv.Elem()
	for i := 0; i < pve.NumField(); i++ {
		fv := pve.Field(i)
		ft := pve.Type().Field(i)
		// ylog.Debug(ft.Name, fv.IsValid(), fv.CanSet())
		if !fv.IsValid() || !fv.CanSet() {
			continue
		}
		switch ft.Type.Kind() {
		case reflect.Map:
			fv.Set(reflect.MakeMap(ft.Type))
		case reflect.Slice:
			fv.Set(reflect.MakeSlice(ft.Type, 0, 0))
		case reflect.Chan:
			fv.Set(reflect.MakeChan(ft.Type, 0))
		case reflect.Struct:
			srv := newValueStruct(ft.Type)
			if !IsNilOrInvalidValue(srv) {
				fv.Set(srv)
			}
		case reflect.Ptr:
			ylog.Debug(ft.Name)
			srv := newValue(ft.Type, nil)
			if !IsNilOrInvalidValue(srv) {
				fv.Set(srv.Elem())
			}
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			fv.SetInt(0)
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
			fv.SetUint(0)
		case reflect.Float32, reflect.Float64:
			fv.SetFloat(0)
		case reflect.Bool:
			fv.SetBool(false)
		case reflect.String:
			fv.SetString("")
		default:
		}
	}
	return pv
}

func newValueMap(t reflect.Type) reflect.Value {
	if t.Kind() == reflect.Ptr {
		pv := reflect.New(t)
		cv := newValueMap(t.Elem())
		pv.Elem().Set(cv)
		return pv
	}
	return reflect.MakeMap(t)
}

func newValueSlice(t reflect.Type) reflect.Value {
	if t.Kind() == reflect.Ptr {
		pv := reflect.New(t)
		cv := newValueSlice(t.Elem())
		pv.Elem().Set(cv)
		return pv
	}
	return reflect.MakeSlice(t, 0, 0)
}


func newValueChan(t reflect.Type) reflect.Value {
	if t.Kind() == reflect.Ptr {
		pv := reflect.New(t)
		cv := newValueChan(t.Elem())
		pv.Elem().Set(cv)
		return pv
	}
	return reflect.MakeChan(t, 0)
}

func newValueScalar(t reflect.Type) reflect.Value {
	if t.Kind() == reflect.Ptr {
		cv := newValueScalar(t.Elem())
		return newPtrOfValue(cv)
	}
	pv := reflect.New(t)
	pve := pv.Elem()
	switch pve.Type().Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		pve.SetInt(0)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		pve.SetUint(0)
	case reflect.Float32, reflect.Float64:
		pve.SetFloat(0)
	case reflect.Bool:
		pve.SetBool(false)
	case reflect.String:
		pve.SetString("")
	default:
	}
	return pve
}


// return new ptr variable of the t type. e.g. T ==> *T
// key is just used for map type data.
// value is used for the base value of the created variable.
func newValue(t reflect.Type, value interface{}) reflect.Value {
	if t == reflect.TypeOf(nil) {
		return reflect.Value{}
	}
	pt := GetPureType(t)
	if IsTypeStruct(pt) {
		return newValueStruct(t)
	} else if IsTypeMap(pt) {
		return newValueMap(t)
	} else if IsTypeSlice(pt) {
		nv := newValueSlice(t)
		if value != nil || !IsNilOrInvalidValue(reflect.ValueOf(value)) {
			nv = setSliceValue(nv, value)
		}
		return nv
	} else if IsTypeChan(pt) {
		return newValueChan(t)
	} else {
		v := newValueScalar(t)
		if IsValueNil(value) {
			return v
		}
		setValueScalar(v, value)
		return v
	}
}

// ptr wraps the given value with pointer: V => *V, *V => **V, etc.
func newPtrOfValue(v reflect.Value) reflect.Value {
	pt := reflect.PtrTo(v.Type()) // create a *T type.
	pv := reflect.New(pt.Elem())  // create a reflect.Value of type *T.
	pv.Elem().Set(v)              // sets pv to point to underlying value of v.
	return pv
}

// NewValue creates an variable of the value based on the type of the typRef.
func NewValue(typRef interface{}, value interface{}) interface{} {
	// v := reflect.ValueOf(value)
	// if IsNilOrInvalidValue(v) {
	// 	return nil
	// }
	nv := newValue(reflect.TypeOf(typRef), value)
	if nv.IsValid() {
		return nv.Interface()
	}
	return nil
}

// SetScalarValue sets the value based on source.
func SetScalarValue(target interface{}, value interface{}) error {
	v := reflect.ValueOf(target)
	if IsNilOrInvalidValue(v) {
		return fmt.Errorf("nil or invalid:: %s", PrintValue(target))
	}
	if !IsValuePtr(v) {
		return fmt.Errorf("no ptr :: %s", PrintValue(target))
	}
	if IsValueStruct(v) || IsValueMap(v) || IsValueSlice(v) {
		return fmt.Errorf("no scalar value :: %s", PrintValue(target))
	}
	return setValueScalar(v, value)
}

// SetMapValue sets the value based on source.
func SetMapValue(target interface{}, key interface{}, value interface{}) error {
	t := reflect.TypeOf(target)
	v := reflect.ValueOf(target)
	if IsValueNil(target) {
		return fmt.Errorf("nil map :: %s", PrintValue(target))
	}
	ylog.Debug("T", t)
	if IsTypeDeep(t, reflect.Map) {
		return setMapValue(v, key, value)
	}
	return fmt.Errorf("no map :: %s", PrintValue(target))
}

// SetSliceValue sets the value based on source.
func SetSliceValue(target interface{}, value interface{}) interface{} {
	t := reflect.TypeOf(target)
	v := reflect.ValueOf(target)
	if IsValueNil(target) {
		return nil
	}
	if IsTypeDeep(t, reflect.Slice) {
		nv := setSliceValue(v, value)
		return nv.Interface()
	}
	return nil
}

