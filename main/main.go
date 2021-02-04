package main

import (
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
}
