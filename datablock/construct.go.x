// Copyright 2017 Google Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Binary uncompressed is an example package showing the usage of ygot for
// an uncompressed schema.
package main

import (
	"fmt"
	"hfr/modeldata/gostruct"
	"hfr/modeldata/gostruct/demo"

	"github.com/openconfig/goyang/pkg/yang"
	"github.com/openconfig/ygot/ygot"
	"github.com/openconfig/ygot/ytypes"
)

// Generate rule to create the example structs:
//go:generate go run ../../github.com/openconfig/ygot/generator/generator.go -path=yang -output_file=gostruct/demo/model.go -package_name=demo -generate_fakeroot -fakeroot_name=device yang/example.yang

var (
	// SchemaTree rearranged by name
	SchemaTree map[string][]*yang.Entry
	Schema *ytypes.Schema
	SchemaRoot *yang.Entry
	DataRoot interface{}
)

func init() {
	s, err := demo.Schema()
	if err != nil {
		panic("Failed to load Schema")
	}
	Schema = s
	SchemaTree = make(map[string][]*yang.Entry)
	for _, val := range demo.SchemaTree {
		entries, _ := SchemaTree[val.Name]
		entries = append(entries, val)
		SchemaTree[val.Name] = entries
		if val.Annotation["schemapath"] == "/" {
			SchemaRoot = val
		}
		fmt.Println(val)
	}
	if SchemaRoot == nil {
		panic("Faled to load SchemaRoot")
	}
}

func find(entry *yang.Entry, keys ...string) *yang.Entry {
	var found *yang.Entry
	if entry == nil {
		return nil
	}
	if len(keys) > 1 {
		found = entry.Dir[keys[0]]
		if found == nil {
			return nil
		}
		found = find(found, keys[1:]...)
	} else {
		found = entry.Dir[keys[0]]
	}
	return found
}



func main() {
	device := Schema.Root
	fmt.Println(device)
	
	// found := find(SchemaRoot, "country")
	// fmt.Println(found)

	e, err := BuildDemo()
	if err != nil {
		panic(err)
	}

	ij, err := DemoInternalJSON(e)
	if err != nil {
		panic(fmt.Sprintf("Internal error: %v", err))
	}
	fmt.Println(ij)

	rj, err := DemoRFC7951JSON(e)
	if err != nil {
		panic(fmt.Sprintf("RFC7951 error: %v", err))
	}
	fmt.Println(rj)
	y := gostruct.YPath{}
	y.Set("network:country[name='United Kingdom']", "country-code")
	gostruct.GetNodes(SchemaRoot, e, &y)
}

// BuildDemo populates a demo instance of the uncompressed GoStructs
// for the example.yang module.
func BuildDemo() (*demo.Device, error) {
	d := &demo.Device{
		Person: ygot.String("robjs"),
	}
	uk, err := d.NewCountry("United Kingdom")
	if err != nil {
		return nil, err
	}
	uk.CountryCode = ygot.String("GB")
	uk.DialCode = ygot.Uint32(44)

	c2, err := d.NewOperator(29636)
	if err != nil {
		return nil, err
	}
	c2.Name = ygot.String("Catalyst2")

	if err := d.Validate(); err != nil {
		return nil, err
	}

	return d, nil
}

// DemoInternalJSON returns internal format JSON for the input
// ucompressed root struct d.
func DemoInternalJSON(d *demo.Device) (string, error) {
	json, err := ygot.EmitJSON(d, &ygot.EmitJSONConfig{
		Format: ygot.Internal,
		Indent: "  ",
	})
	if err != nil {
		return "", err
	}
	return json, nil
}

// DemoRFC7951JSON returns RFC7951 JSON for the input uncompressed
// root struct d.
func DemoRFC7951JSON(d *demo.Device) (string, error) {
	json, err := ygot.EmitJSON(d, &ygot.EmitJSONConfig{
		Format: ygot.RFC7951,
		Indent: "  ",
		RFC7951Config: &ygot.RFC7951JSONConfig{
			AppendModuleName: true,
		},
	})
	if err != nil {
		return "", err
	}
	return json, nil

}
