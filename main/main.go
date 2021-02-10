package main

import (
	"fmt"
	"git.ash.lt/allrss/cssparser/parser"
	"io/ioutil"
	"log"
)

func main() {
	bs, err := ioutil.ReadFile("./test2.css")
	if err != nil {
		log.Fatal("error while read test.css ", err)
	}

	style, err := parser.Parse(string(bs))

	if err != nil {
		log.Fatalln(err)
	}

	style.FindByKey("color", func(rule *parser.Rule) {
		//rule.Value = "blabla"
		//rule.Delete()
	})

	fmt.Println(style.StringCSS())

	//fmt.Println(style)

	//s, _ := json.MarshalIndent(style.Model, "", "\t")
	//fmt.Println(string(s))
}
