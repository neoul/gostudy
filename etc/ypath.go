package main

import (
	"fmt"
	"strings"
)

type YPathElem struct {
	Prefix string
	Name string
	Key map[string]string
}

func (elem *YPathElem) GetKey() map[string]string {
	if elem != nil {
		return elem.Key
	}
	return nil
}

func (elem *YPathElem) GetName() string {
	if elem != nil {
		return elem.Name
	}
	return ""
}

func (elem *YPathElem) GetPrefix() string {
	if elem != nil {
		return elem.Prefix
	}
	return ""
}

func (elem *YPathElem) GetPrefixAndName() string {
	if elem == nil {
		return ""
	}
	if elem.Prefix != "" {
		return fmt.Sprintf("%s:%s", elem.Prefix, elem.Name)
	} else {
		return elem.Name
	}
}

type YPath struct {
	Elem []*YPathElem
}

func makeYPathElem(s string) *YPathElem {
	var prefix string
	prefixend := strings.Index(s, ":")
	if prefixend >= 0 {
		prefix = s[:prefixend]
		s = s[prefixend+1:]
	}
	keystart := strings.Index(s, "[")
	keyend := strings.LastIndex(s, "]")
	if keystart >= 0 && keyend >= 0 && keystart < keyend {
		elemName := s[:keystart]
		keyNameValue := s[keystart+1:keyend]
		keyvalstart := strings.Index(keyNameValue, "=")
		if keyvalstart >= 0 {
			keyname := keyNameValue[:keyvalstart]
			keyvalue := keyNameValue[keyvalstart+1:]
			keyvalue = strings.Trim(keyvalue, " '\"")
			elem := &(YPathElem{Name: elemName, Prefix: prefix, Key: make(map[string]string)})
			elem.Key[keyname] = keyvalue
			return elem
		}
	}
	return &(YPathElem{Name: s, Prefix: prefix})
}

func (ypath *YPath) Set(pathnodes ...string) {
	if ypath.Elem == nil {
		ypath.Elem = make([]*YPathElem, 0, 12)
	}
	for _, node := range pathnodes {
		elem := makeYPathElem(node)

		ypath.Elem = append(ypath.Elem, elem)
	}
}

func (ypath *YPath) Print() {
	for _, elem := range ypath.Elem {
		fmt.Printf("%s", elem.GetPrefixAndName())
		if elem.Key != nil || len(elem.Key) > 0 {
			fmt.Printf(" %s", elem.GetKey())
		}
		fmt.Printf("\n")
	}
}

func (ypath *YPath) GetElem() []*YPathElem {
	return ypath.Elem
}

// func main() {
// 	y := &(YPath{})
// 	y.Set("openconfig-interfaces:interfaces", "interface[name=1/1]")
// 	y.Print()
// }
