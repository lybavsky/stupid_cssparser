package main

import (
	"fmt"
	"io/ioutil"
	"log"
)

func main() {
	bs, err := ioutil.ReadFile("./test.css")
	if err != nil {
		log.Fatal("error while read test.css ", err)
	}

	//log.Println(string(bs))

	runes := []rune(string(bs))

	var inside_comment bool = false
	var inside_string rune = 0
	var inside_bracket int = 0

	//blocks := []string{}

	var curr_block []rune = []rune{}
	var comment_block []rune = make([]rune, 0)

	//status := func(r rune) {
	//
	//	fmt.Println("rune: ", string(r), "\n\tcomment: ", inside_comment,
	//		", string: ", inside_string,
	//		"bra: ", inside_bracket,
	//		"\n\tbuf: ", string(curr_block),
	//		"\n\tcbuf: ", string(comment_block))
	//}

	for i := 0; i < len(runes); i++ {
		cur_rune := runes[i]
		//status(cur_rune)

		//Если просто ползем по строке и натыкаемся на начало коммента..
		if !inside_comment && inside_string == 0 && cur_rune == '/' && len(runes) > i && runes[i+1] == '*' {
			inside_comment = true
			i++
			continue
		} else {
			//Если уже внутри коммента
			if inside_comment {
				//Если закрывающий символ
				if cur_rune == '*' && len(runes) > i && runes[i+1] == '/' {
					inside_comment = false
					i++
					log.Println("Found comment: " + string(comment_block))
					comment_block = []rune{}
				} else {
					comment_block = append(comment_block, cur_rune)
				}
				continue

			}

		}

		//Если не внутри строки, но попадаем на символ открывающей кавычки, считаем, что попали внутрь строки
		if inside_string == 0 && (cur_rune == '\'' || cur_rune == '"') {
			inside_string = cur_rune
		} else {
			//Если попали на соответствующую закрывающую кавычку
			if inside_string != 0 && inside_string == cur_rune {
				inside_string = 0
			}
		}

		//Если попали на пустой символ, в нынешнем блоке последний символ не пробел
		if inside_string == 0 && cur_rune == '\n' || cur_rune == '\r' {
			if len(curr_block) > 0 && curr_block[len(curr_block)-1] != ' ' {
				curr_block = append(curr_block, ' ')
			}
			continue
		}

		//Если попали на второй подряд пробел, пропускаем
		if inside_string == 0 && cur_rune == ' ' && curr_block[len(curr_block)-1] == ' ' {
			continue
		}

		//Добавлям наш символ к блоку
		curr_block = append(curr_block, cur_rune)

		//Если попадаем на открывающую скобку, считаем, что провалились в блок
		if inside_string == 0 && cur_rune == '{' {
			inside_bracket++
		}

		//Если попали на последнюю закрыващую скобку, считаем, что блок закончился
		if inside_string == 0 && cur_rune == '}' {
			inside_bracket--
			if inside_bracket < 0 {
				log.Fatalln("Failed parse css - negative number of strings")
			}

			if inside_bracket == 0 {
				fmt.Println("Block: ", string(curr_block))
				curr_block = []rune{}
				continue
			}
		}

		//Если не внутри скобок и поймали точку с запятой
		if inside_bracket == 0 && cur_rune == ';' {
			fmt.Println("Block: ", string(curr_block))
			curr_block = []rune{}
			continue
		}

	//TODO: Дальше, видимо, парсить непосредственно блоки
	}
}

type BlockType int

const (
	AtType      BlockType = 0
	CommentType           = iota
	RuleType
)

type Media struct {
	Selector string
	Rules    []Rule
}
type Rule struct {
	Selector string
	Value    string
}
