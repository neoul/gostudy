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
package main // import "github.com/neoul/gostudy/datablock/ydb2ygot"

import (
	"fmt"

	"github.com/neoul/gostudy/datablock/model/object"

	"github.com/openconfig/goyang/pkg/yang"
	"github.com/openconfig/ygot/ytypes"
)

// Generate rule to create the example structs:
//go:generate go run ../../../github.com/openconfig/ygot/generator/generator.go -path=yang -output_file=object/example.go -package_name=object -generate_fakeroot -fakeroot_name=device yang/example.yang

var (
	// SchemaTree rearranged by name
	SchemaTree map[string][]*yang.Entry
	Schema *ytypes.Schema
	SchemaRoot *yang.Entry
	DataRoot interface{}
)

func init() {
	schema, err := object.Schema()
	if err != nil {
		panic("Failed to load Schema")
	}
	Schema = schema
	SchemaTree = make(map[string][]*yang.Entry)
	for _, branch := range object.SchemaTree {
		entries, _ := SchemaTree[branch.Name]
		entries = append(entries, branch)
		for _, leaf := range branch.Dir {
			entries = append(entries, leaf)
		}
		SchemaTree[branch.Name] = entries
		if branch.Annotation["schemapath"] == "/" {
			SchemaRoot = branch
		}
		fmt.Println(branch)
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
}
