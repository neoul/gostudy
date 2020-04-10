package main

import (
	"encoding/json"
	"log"
)

type DataBlock struct {
}

func main() {
	doc := `
	{
		"name": "maria",
		"age": 10,
		"info": {
			"level": 4,
			"dex": 10
		}
	}`

	var data map[interface{}]interface{}
	err := json.Unmarshal([]byte(doc), &data)
	log.Println(err)
	log.Println(data)

}
