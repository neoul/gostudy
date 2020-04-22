package reflection // import "github.com/neoul/gostudy/datablock/reflection"

import (
	"fmt"
	"os"
	"reflect"
	"strconv"

	"github.com/neoul/gostudy/datablock/log"
)

var (
	ylog log.Log
)

func init() {
	ylog = log.NewLog("reflection", os.Stdout)
}

func typeString(t reflect.Type) string {
	return fmt.Sprintf("%s", t)
}

// getBaseType returns not reflect.Ptr type.
func getBaseType(t reflect.Type) reflect.Type {
	for ; t.Kind() == reflect.Ptr; t = t.Elem() {
	}
	return t
}

// isTypeDeep reports whether t is k type.
func isTypeDeep(t reflect.Type, kinds ...reflect.Kind) bool {
	for ; t.Kind() == reflect.Ptr; t = t.Elem() {
	}
	for _, k := range kinds {
		if t.Kind() == k {
			return true
		}
	}
	return false
}

// isReferenceType returns true if t is a map, slice or channel
func isReferenceType(t reflect.Type) bool {
	switch t.Kind() {
	case reflect.Slice, reflect.Chan, reflect.Map:
		return true
	}
	return false
}

// isTypeStruct reports whether t is a struct type.
func isTypeStruct(t reflect.Type) bool {
	return t.Kind() == reflect.Struct
}

// isTypeStructPtr reports whether v is a struct ptr type.
func isTypeStructPtr(t reflect.Type) bool {
	if t == reflect.TypeOf(nil) {
		return false
	}
	return t.Kind() == reflect.Ptr && t.Elem().Kind() == reflect.Struct
}

// isTypeSlice reports whether v is a slice type.
func isTypeSlice(t reflect.Type) bool {
	return t.Kind() == reflect.Slice
}

// isTypeSlicePtr reports whether v is a slice ptr type.
func isTypeSlicePtr(t reflect.Type) bool {
	if t == reflect.TypeOf(nil) {
		return false
	}
	return t.Kind() == reflect.Ptr && t.Elem().Kind() == reflect.Slice
}

// isTypeMap reports whether v is a map type.
func isTypeMap(t reflect.Type) bool {
	if t == reflect.TypeOf(nil) {
		return false
	}
	return t.Kind() == reflect.Map
}

// isTypeInterface reports whether v is an interface.
func isTypeInterface(t reflect.Type) bool {
	if t == reflect.TypeOf(nil) {
		return false
	}
	return t.Kind() == reflect.Interface
}

// isTypeSliceOfInterface reports whether v is a slice of interface.
func isTypeSliceOfInterface(t reflect.Type) bool {
	if t == reflect.TypeOf(nil) {
		return false
	}
	return t.Kind() == reflect.Slice && t.Elem().Kind() == reflect.Interface
}

// isTypeChan reports whether v is a slice type.
func isTypeChan(t reflect.Type) bool {
	return t.Kind() == reflect.Chan
}

// AreSameType returns true if t1 and t2 has the same reflect.Type,
// otherwise it returns false.
func AreSameType(t1 reflect.Type, t2 reflect.Type) bool {
	b1 := getBaseType(t1)
	b2 := getBaseType(t2)
	return b1 == b2
}

// isNilOrInvalidValue reports whether v is nil or reflect.Zero.
func isNilOrInvalidValue(v reflect.Value) bool {
	return !v.IsValid() || (v.Kind() == reflect.Ptr && v.IsNil()) || isValueNil(v.Interface())
}

// isValueNil returns true if either value is nil, or has dynamic type {ptr,
// map, slice} with value nil.
func isValueNil(value interface{}) bool {
	if value == nil {
		return true
	}
	switch reflect.TypeOf(value).Kind() {
	case reflect.Slice, reflect.Ptr, reflect.Map:
		return reflect.ValueOf(value).IsNil()
	}
	return false
}

// isValueNilOrDefault returns true if either isValueNil(value) or the default
// value for the type.
func isValueNilOrDefault(value interface{}) bool {
	if isValueNil(value) {
		return true
	}
	if !isValueScalar(reflect.ValueOf(value)) {
		// Default value is nil for non-scalar types.
		return false
	}
	return value == reflect.New(reflect.TypeOf(value)).Elem().Interface()
}

// isValuePtr reports whether v is a ptr type.
func isValuePtr(v reflect.Value) bool {
	return v.Kind() == reflect.Ptr
}

// isValueInterface reports whether v is an interface type.
func isValueInterface(v reflect.Value) bool {
	return v.Kind() == reflect.Interface
}

// isValueStruct reports whether v is a struct type.
func isValueStruct(v reflect.Value) bool {
	return v.Kind() == reflect.Struct
}

// isValueStructPtr reports whether v is a struct ptr type.
func isValueStructPtr(v reflect.Value) bool {
	return v.Kind() == reflect.Ptr && isValueStruct(v.Elem())
}

// isValueMap reports whether v is a map type.
func isValueMap(v reflect.Value) bool {
	return v.Kind() == reflect.Map
}

// isValueSlice reports whether v is a slice type.
func isValueSlice(v reflect.Value) bool {
	return v.Kind() == reflect.Slice
}

// isValueScalar reports whether v is a scalar type.
func isValueScalar(v reflect.Value) bool {
	if isNilOrInvalidValue(v) {
		return false
	}
	if isValuePtr(v) {
		if v.IsNil() {
			return false
		}
		v = v.Elem()
	}
	return !isValueStruct(v) && !isValueMap(v) && !isValueSlice(v)
}


// isValueInterfaceToStructPtr reports whether v is an interface that contains a
// pointer to a struct.
func isValueInterfaceToStructPtr(v reflect.Value) bool {
	return isValueInterface(v) && isValueStructPtr(v.Elem())
}

// isStructValueWithNFields returns true if the reflect.Value representing a
// struct v has n fields.
func isStructValueWithNFields(v reflect.Value, n int) bool {
	return isValueStruct(v) && v.NumField() == n
}

// isSimpleType - true if built-in simple variable type
func isSimpleType(t reflect.Type) bool {
	switch t.Kind() {
	case reflect.Ptr:
		return isSimpleType(t.Elem())
	case reflect.Array, reflect.Chan, reflect.Func, reflect.Interface,
		reflect.Map, reflect.Slice, reflect.Struct,
		reflect.UnsafePointer, reflect.Complex64, reflect.Complex128:
		return false
	default:
		return true
	}
}

// isValidDeep reports whether v is valid.
func isValidDeep(v reflect.Value) bool {
	for ; v.Kind() == reflect.Ptr; v = v.Elem() {
		if !v.IsValid() {
			return false
		}
	}
	if !v.IsValid() {
		return false
	}
	return true
}

// isNilDeep reports whether v is nil
func isNilDeep(v reflect.Value) bool {
	for ; v.Kind() == reflect.Ptr; v = v.Elem() {
		if v.IsNil() {
			return true
		}
	}
	return false
}

var maxValueStringLen = 150

// ValueString returns a string representation of value which may be a value, ptr,
// or struct type.
func ValueString(value interface{}) string {
	v := reflect.ValueOf(value)
	out := valueString(v, 0)
	if len(out) > maxValueStringLen {
		out = out[:maxValueStringLen] + "..."
	}
	return out
}

// ValueStrInternal is the internal implementation of ValueString.
func valueString(v reflect.Value, ptrcnt int) string {
	var out string
	if isNilOrInvalidValue(v) {
		return fmt.Sprintf("%s{?}", v.Type())
	}
	// ylog.Debug("v:", v)
	switch v.Kind() {
	case reflect.Ptr:
		ptrcnt++
		out = "*" + valueString(v.Elem(), ptrcnt)
	case reflect.Slice:
		out = fmt.Sprintf("%s{", v.Type())
		for i := 0; i < v.Len(); i++ {
			if i != 0 {
				out += ","
			}
			out += ValueString(v.Index(i).Interface())
		}
		out += "}"
	case reflect.Struct:
		comma := false
		t := v.Type()
		out = fmt.Sprintf("%s{", v.Type())
		for i := 0; i < v.NumField(); i++ {
			fv := v.Field(i)
			ft := t.Field(i)
			ylog.Debug(ft.Name, v.Type(), ft.Type)
			if AreSameType(ft.Type, t) {
				continue
			}
			if comma {
				out += ","
			}
			if fv.CanInterface() {
				out += fmt.Sprintf("%s:%v", ft.Name, ValueString(fv.Interface()))
			} else {
				out += fmt.Sprintf("%s:%v", ft.Name, fv)
			}
			comma = true
		}
		out += "}"
	case reflect.Map:
		comma := false
		out = fmt.Sprintf("%s{", v.Type())
		iter := v.MapRange()
		for iter.Next() {
			k := iter.Key()
			e := iter.Value()
			if comma {
				out += ","
			}
			out += fmt.Sprintf("%v:%s", k, ValueString(e.Interface()))
			comma = true
		}
		out += "}"
	default:
		out = fmt.Sprintf("%s{", v.Type())
		for i := 0; i < ptrcnt; i++ {
			out = out + "&"
		}
		out = out + fmt.Sprintf("%v}", v)
	}
	if len(out) > maxValueStringLen {
		out = out[:maxValueStringLen] + "..."
	}
	return out
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
		if st.Elem().Kind() == reflect.Ptr {
			nv := setSliceValue(sv.Elem(), element)
			sv.Elem().Set(nv)
			return sv
		} 
		if st.Elem().Kind() == reflect.Slice {
			et := st.Elem().Elem()
			ev := newValue(et, element)
			sv.Elem().Set(reflect.Append(sv.Elem(), ev))
			// ylog.Info("sv Slice:::::", sv)
			return sv
		}
	}
	ev := newValue(st.Elem(), element)
	num := sv.Len()
	nslice := reflect.MakeSlice(st, num+1, num+1)
	reflect.Copy(nslice, sv)
	nslice.Index(num).Set(ev)
	return nslice
}

func copySliceValue(v reflect.Value) reflect.Value {
	sv := v
	st := sv.Type()
	if st.Kind() == reflect.Ptr {
		if st.Elem().Kind() == reflect.Ptr {
			cv := copySliceValue(sv.Elem())
			pv := newPtrOfValue(cv)
			return pv
		}
		sv = sv.Elem()
		st = st.Elem()
	}
	num := sv.Len()
	nslice := reflect.MakeSlice(st, num, num)
	reflect.Copy(nslice, sv)
	if sv != v {
		return newPtrOfValue(nslice)
	}
	return nslice
}

func setValueScalar(v reflect.Value, value interface{}) error {
	dv := v
	if dv.Kind() == reflect.Ptr {
		dv = v.Elem()
		if dv.Kind() == reflect.Ptr {
			if dv.IsNil() { // e.g. **type
				dv = reflect.New(dv.Type().Elem())
				// ylog.Debug(ValueString(dv.Interface()))
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
					return fmt.Errorf("overflowInt: %s", ValueString(val))
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
					return fmt.Errorf("OverflowUint: %s", ValueString(val))
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
					return fmt.Errorf("OverflowFloat: %s", ValueString(val))
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
	return fmt.Errorf("Not Convertible: %s", ValueString(v.Interface()))
}


func getStructFieldName(ft reflect.StructField) (string, string) {
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
	return prefix, name
}

func searchStructField(st reflect.Type, sv reflect.Value, name string) (reflect.StructField, reflect.Value, bool) {
	var fv reflect.Value
	ft, ok := st.FieldByName(name)
	if ok {
		fv = sv.FieldByName(name)
		return ft, fv, true
	}

	for i := 0; i < sv.NumField(); i++ {
		fv := sv.Field(i)
		ft := st.Field(i)
		if !fv.IsValid() || !fv.CanSet() {
			continue
		}
		_, name := getStructFieldName(ft)
		if name != "" {
			return ft, fv, true
		}
	}
	return ft, reflect.Value{}, false
}

func setStructField(sv reflect.Value, fieldName interface{}, fieldValue interface{}) error {
	if sv.Kind() == reflect.Ptr {
		if sv.IsNil() {
			cv := newValueStruct(sv.Type())
			sv.Set(cv)
		} else {
			return setStructField(sv.Elem(), fieldName, fieldValue)
		}
	}
	fieldname := reflect.TypeOf("")
	kv := newValue(fieldname, fieldName)
	ft, fv, ok := searchStructField(sv.Type(), sv, kv.Interface().(string))
	ylog.Debug(ft, fv, ok)
	if !ok {
		return fmt.Errorf("not found %s.%s", sv.Type(), fieldname)
	}
	nv := NewValue(ft.Type, fieldValue)
	fv.Set(nv)
	return nil
}

func newValueStruct(t reflect.Type) reflect.Value {
	if t.Kind() == reflect.Ptr {
		cv := newValueStruct(t.Elem())
		pv := newPtrOfValue(cv)
		return pv
	}

	pv := reflect.New(t)
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
			if !isNilOrInvalidValue(srv) {
				fv.Set(srv)
			}
		case reflect.Ptr:
			// ylog.Debug(ft.Name, ft.Type)
			srv := newValue(ft.Type, nil)
			// ylog.Debug(fv, srv)
			if !isNilOrInvalidValue(srv) {
				fv.Set(srv)
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
		cv := newValueMap(t.Elem())
		pv := newPtrOfValue(cv)
		return pv
	}
	return reflect.MakeMap(t)
}

func newValueSlice(t reflect.Type) reflect.Value {
	if t.Kind() == reflect.Ptr {
		cv := newValueSlice(t.Elem())
		pv := newPtrOfValue(cv)
		return pv
	}
	return reflect.MakeSlice(t, 0, 0)
}

func newValueChan(t reflect.Type) reflect.Value {
	if t.Kind() == reflect.Ptr {
		cv := newValueChan(t.Elem())
		pv := newPtrOfValue(cv)
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
	pt := getBaseType(t)
	if isTypeStruct(pt) {
		return newValueStruct(t)
	} else if isTypeMap(pt) {
		return newValueMap(t)
	} else if isTypeSlice(pt) {
		nv := newValueSlice(t)
		if value != nil || !isNilOrInvalidValue(reflect.ValueOf(value)) {
			nv = setSliceValue(nv, value)
		}
		return nv
	} else if isTypeChan(pt) {
		return newValueChan(t)
	} else {
		v := newValueScalar(t)
		if isValueNil(value) {
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

// NewValue - returns new Value based on t type
func NewValue(t reflect.Type, values ...interface{}) reflect.Value {
	if t == reflect.TypeOf(nil) {
		return reflect.Value{}
	}
	pt := getBaseType(t)
	switch pt.Kind() {
	case reflect.Array, reflect.Complex64, reflect.Complex128, reflect.Chan:
		return reflect.Value{}
	case reflect.Struct:
		var key interface{} = nil
		nv := newValueStruct(t)
		for _, val := range values {
			if key == nil {
				key = val
			} else {
				ylog.Info(key, val)
				setStructField(nv, key, val)
				key = nil
			}
		}
		if key != nil {
			ylog.Warningf("Single key value ignored '%v:none'", key)
		}
		return nv
	case reflect.Map:
		var key interface{} = nil
		nv := newValueMap(t)
		for _, val := range values {
			if key == nil {
				key = val
			} else {
				setMapValue(nv, key, val)
				key = nil
			}
		}
		if key != nil {
			ylog.Warningf("Single key value ignored '%v:none'", key)
		}
		return nv
	case reflect.Slice:
		nv := newValueSlice(t)
		for _, val := range values {
			nv = setSliceValue(nv, val)
		}
		return nv
	default:
		nv := newValueScalar(t)
		for _, val := range values {
			err:= setValueScalar(nv, val)
			if err != nil {
				ylog.Warningf("Not settable value inserted '%s'", ValueString(val))
			}
		}
		return nv
	}
}


// SetValue - returns new Value based on t type
func SetValue(v reflect.Value, values ...interface{}) reflect.Value {
	if !v.IsValid() {
		return v
	}

	t := v.Type()
	pt := getBaseType(t)
	switch pt.Kind() {
	case reflect.Array, reflect.Complex64, reflect.Complex128, reflect.Chan:
		return v
	case reflect.Struct:
		var key interface{} = nil
		for _, val := range values {
			if key == nil {
				key = val
			} else {
				setStructField(v, key, val)
				key = nil
			}
		}
		return v
	case reflect.Map:
		var key interface{} = nil
		for _, val := range values {
			if key == nil {
				key = val
			} else {
				setMapValue(v, key, val)
				key = nil
			}
		}
		return v
	case reflect.Slice:
		var nv reflect.Value
		isnil := isNilDeep(v)
		if isnil {
			nv = newValueSlice(t)
		} else {
			nv = copySliceValue(v)
		}
		for _, val := range values {
			nv = setSliceValue(nv, val)
		}
		if t.Kind() == reflect.Ptr {
			if v.Elem().CanSet() {
				v.Elem().Set(nv.Elem())
			} else {
				ylog.Warningf("Not settable variable '%s'", ValueString(v.Interface()))
			}
		}
		return nv
	default:
		nv := v
		if t.Kind() != reflect.Ptr {
			nv = newValueScalar(t)
		}
		for _, val := range values {
			err:= setValueScalar(nv, val)
			if err != nil {
				ylog.Warningf("Not settable value inserted '%s'", ValueString(val))
			}
		}
		return nv
	}
}
