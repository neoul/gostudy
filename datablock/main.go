package main

import (
	"os"
	"reflect"

	"github.com/neoul/gostudy/datablock/log"
	yref "github.com/neoul/gostudy/datablock/reflection"
)


func main() {
	var out reflect.Value
	var i int = 1
	var ii *int
	var iii **int

	ylog := log.NewLog("main", os.Stdout)
	out = yref.NewValue(reflect.TypeOf(i), 10)
	ylog.Debug(yref.ValueString(out.Interface()))
	out = yref.NewValue(reflect.TypeOf(&i), 10)
	ylog.Debug(yref.ValueString(out.Interface()))

	out = yref.NewValue(reflect.TypeOf(ii), 20)
	ylog.Debug(yref.ValueString(out.Interface()))
	out = yref.NewValue(reflect.TypeOf(iii), 30)
	ylog.Debug(yref.ValueString(out.Interface()))

	var si []int
	var ssi *[]int
	var ssii *[]*int
	out = yref.NewValue(reflect.TypeOf(si), 10, 20, 30)
	si = out.Interface().([]int)
	ylog.Debug(yref.ValueString(si))
	out = yref.NewValue(reflect.TypeOf(ssi), 10, 20, 30)
	ylog.Debug(yref.ValueString(out.Interface()))
	out = yref.NewValue(reflect.TypeOf(ssii), 10, 20, 30)
	ylog.Debug(yref.ValueString(out.Interface()))

	var mi map[int]int
	var mmi *map[int]int
	var mmii *map[int]*int
	out = yref.NewValue(reflect.TypeOf(mi), 10, 20, 30, 40)
	ylog.Debug(yref.ValueString(out.Interface()))
	out = yref.NewValue(reflect.TypeOf(mmi), 10, 20, 30, 40)
	ylog.Debug(yref.ValueString(out.Interface()))
	out = yref.NewValue(reflect.TypeOf(mmii), 10, 20, 30, 40)
	ylog.Debug(yref.ValueString(out.Interface()))

	type structi struct {
		I int
		II *int
		III map[int]int
	}
	var sti structi
	out = yref.NewValue(reflect.TypeOf(sti))
	ylog.Debug(yref.ValueString(out.Interface()))
	ylog.Debug("???")
	iim := make(map[string]string)
	iim["10"] = "11"
	out = yref.NewValue(reflect.TypeOf(sti), "I", "10", "II", 20, "III", iim)
	ylog.Debug("???")
	ylog.Debug(yref.ValueString(out.Interface()))

	ylog.Debug("")
	ylog.Debug("[SetValue(value)]")
	ylog.Debug("")
	
	out = yref.SetValue(reflect.ValueOf(i), 1000)
	ylog.Debug("in ", yref.ValueString(i))
	ylog.Debug("out", yref.ValueString(out.Interface()))

	out = yref.SetValue(reflect.ValueOf(&i), 2000)
	ylog.Debug("in ", yref.ValueString(i))
	ylog.Debug("out", yref.ValueString(out.Interface()))

	ylog.Debug("")
	ylog.Debug("[SetValue(slice)]")
	out = yref.SetValue(reflect.ValueOf(si), 1000)
	ylog.Debug("in : si", yref.ValueString(si))
	ylog.Debug("out: si", yref.ValueString(out.Interface()))
	
	out = yref.SetValue(reflect.ValueOf(&si), 2000)
	ylog.Debug("in : &si", yref.ValueString(&si))
	ylog.Debug("out: &si", yref.ValueString(out.Interface()))
	
	ylog.Debug("")
	// ssi = &si
	out = yref.SetValue(reflect.ValueOf(ssi), 3000)
	ylog.Debug("in : ssi", yref.ValueString(ssi))
	ylog.Debug("out: ssi", yref.ValueString(out.Interface()))

	out = yref.SetValue(reflect.ValueOf(&ssi), 4000)
	ylog.Debug("in : ssi", yref.ValueString(&ssi))
	ylog.Debug("out: ssi", yref.ValueString(out.Interface()))

	out = yref.SetValue(reflect.ValueOf(ssii), 5000, 6000)
	ylog.Debug("in : ssii", yref.ValueString(ssii))
	ylog.Debug("out: ssii", yref.ValueString(out.Interface()))
	
	out = yref.SetValue(reflect.ValueOf(&ssii), 7000, 8000)
	ylog.Debug("in : ssii", yref.ValueString(&ssii))
	ylog.Debug("out: ssii", yref.ValueString(out.Interface()))

	// ssi = &si
	// out = yref.SetValue(reflect.ValueOf(ssi), 3000)
	// ylog.Debug("in : ssi", yref.ValueString(ssi))
	// ylog.Debug("out: ssi", yref.ValueString(out.Interface()))

	// out = yref.SetValue(reflect.ValueOf(&ssi), 4000)
	// ylog.Debug("in : ssi", yref.ValueString(&ssi))
	// ylog.Debug("out: ssi", yref.ValueString(out.Interface()))

	// out = yref.SetValue(reflect.ValueOf(mi), 10, 20, 30, 40)
	// ylog.Debug("in : mi", yref.ValueString(mi))
	// ylog.Debug("out: mi", yref.ValueString(out.Interface()))
}