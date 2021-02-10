package main

import (
	"fmt"
	"git.ash.lt/allrss/cssparser/parser"
	"io/ioutil"
	"log"
	"reflect"
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

	//fmt.Println(style)

	style.FindByKey("color", func(rule *parser.Rule) {
		rule.Value = "black"
		rule.Parent.Rules = []*parser.Rule{}
		fmt.Println(reflect.TypeOf(rule.Parent))
		//rule.Delete()
	})

	//fmt.Println(style.StringCSS())

	//fmt.Println(style)

	//s, _ := json.MarshalIndent(style.Model, "", "\t")
	//fmt.Println(string(s))
}
