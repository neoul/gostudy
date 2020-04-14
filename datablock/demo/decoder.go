package main // import "github.com/neoul/gostudy/datablock/demo"

import (
	"log"
	"os"

	"github.com/neoul/libydb/go/ydb"
)

func main() {
	db, close := ydb.Open("mydb")
	defer close()
	// ydb.SetLog(ydb.LogDebug)

	r, err := os.Open("demo.yaml")
	defer r.Close()
	if err != nil {
		log.Fatalln(err)
	}
	dec := db.NewDecoder(r)
	dec.Decode()

	w, err := os.Create("result.yaml")
	defer w.Close()
	if err != nil {
		log.Fatalln(err)
	}
	enc := db.NewEncoder(w)
	enc.Encode()
}
