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
	out = yref.NewValue(reflect.TypeOf(ii), 20)
	ylog.Debug(yref.ValueString(out.Interface()))
	out = yref.NewValue(reflect.TypeOf(iii), 30)
	ylog.Debug(yref.ValueString(out.Interface()))

	var si []int
	var ssi *[]int
	var ssii *[]*int
	out = yref.NewValue(reflect.TypeOf(si), 10, 20, 30)
	ylog.Debug(yref.ValueString(out.Interface()))
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
	out = yref.NewValue(reflect.TypeOf(sti), "I", "10", "II", 20, "III", "10")
	ylog.Debug(yref.ValueString(out.Interface()))
}