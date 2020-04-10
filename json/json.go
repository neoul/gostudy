package main

import (
	"encoding/json"
	"fmt"
)

func marshal() {
	data := make(map[string]interface{}) // 문자열을 키로하고 모든 자료형을 저장할 수 있는 맵 생성
	data["name"] = "maria"
	data["age"] = 10
	doc, _ := json.Marshal(data) // 맵을 JSON 문서로 변환
	fmt.Println(string(doc)) // {"age":10,"name":"maria"}: 문자열로 변환하여 출력
}

func unmarshal() {
	doc := `
	{
		"name": "maria",
		"age": 10,
		"info": {
			"level": 4,
			"dex": 10
		}
	}`

	var data map[string]interface{}
	err := json.Unmarshal([]byte(doc), &data)
	fmt.Println(err)
	fmt.Println(data)
}


type Author struct {
	Name  string
	Email string
}

type Comment struct {
	Id      uint64
	Author  Author // Author 구조체
	Content string
}

type Article struct {
	Id         uint64
	Title      string
	Author     Author    // Author 구조체
	Content    string
	Recommends []string  // 문자열 배열
	Comments   []Comment // Comment 구조체 배열
}

func main() {
	unmarshal()
	marshal()
	doc := `
	[{
		"Id": 1,
		"Title": "Hello, world!",
		"Author": {
			"Name": "Maria",
			"Email": "maria@example.com"
		},
		"Content": "Hello~",
		"Recommends": [
			"John",
			"Andrew"
		],
		"Comments": [{
			"id": 1,
			"Author": {
				"Name": "Andrew",
				"Email": "andrew@hello.com"
			},
			"Content": "Hello Maria",
			"A": "TEST"
		}]
	}]
	`

	var data []Article // JSON 문서의 데이터를 저장할 구조체 슬라이스 선언

	json.Unmarshal([]byte(doc), &data) // doc의 내용을 변환하여 data에 저장

	fmt.Println(data) // [{1 Hello, world! {Maria maria@exa... (생략)
}

