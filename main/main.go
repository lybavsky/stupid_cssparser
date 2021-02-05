package main

import (
	"encoding/json"
	"fmt"
	"git.ash.lt/allrss/cssparser/parser"
	"io/ioutil"
	"log"
)

func main() {
	bs, err := ioutil.ReadFile("./test.css")
	if err != nil {
		log.Fatal("error while read test.css ", err)
	}

	style, err := parser.Parse(string(bs))

	if err != nil {
		log.Fatalln(err)
	}

	log.Println(style)

	s, _ := json.MarshalIndent(style, "", "\t")
	fmt.Println(string(s))
}
