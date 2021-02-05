package parser

import (
	"errors"
	"fmt"
	"log"
	"strings"
)

func Parse(string_css string) (stylesheet StyleSheet, err error) {
	//Сразу сделаем trim
	string_css=strings.Trim(string_css," \r\n")

	fmt.Println("Will parse: ",string_css)

	atsInherited := make([]StyleSheet, 0)
	blocks := make([]Block, 0)
	ats := make([]string, 0)
	imports := make([]string, 0)

	runes := []rune(string_css)

	var inside_comment bool = false
	var inside_string rune = 0
	var inside_bracket int = 0

	var curr_block []rune = []rune{}

	var comment_block []rune = make([]rune, 0)

	for i := 0; i < len(runes); i++ {
		cur_rune := runes[i]

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
					//fmt.Println("Found comment: " + string(comment_block))
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
		if inside_string == 0 && cur_rune == ' ' && (len(curr_block)==0 || curr_block[len(curr_block)-1] == ' ' ) {
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
				return StyleSheet{}, errors.New("Failed parse css - negative number of strings")
			}

		}

		if (inside_string == 0 && cur_rune == '}' && inside_bracket == 0) ||
			(inside_bracket == 0 && cur_rune == ';') {
			//fmt.Println("Block: ", string(curr_block))

			curr_block = []rune(strings.Trim(string(curr_block), " "))

			if curr_block[0] == '@' {
				if strings.HasPrefix(string(curr_block), "@media") ||
					strings.HasPrefix(string(curr_block), "@supports")  ||
					strings.HasPrefix(string(curr_block), "@document")  {
					atInherited, err := parseAtInherited(curr_block)
					if err != nil {
						return StyleSheet{}, errors.New("Error while parse atInherited " + string(curr_block) + ": " + err.Error())
					}
					atsInherited = append(atsInherited, atInherited)
				} else if strings.HasPrefix(string(curr_block), "@import") {
					imp, err := parseImport(curr_block)
					if err != nil {
						return StyleSheet{}, errors.New("Error while parse import " + string(curr_block) + ": " + err.Error())
					}
					imports = append(imports, imp)
				} else {
					at, err := parseAt(curr_block)
					if err != nil {
						return StyleSheet{}, errors.New("Error while parse at " + string(curr_block) + ": " + err.Error())
					}
					ats = append(ats, at)
				}
			} else {
				block, err := parseBlock(curr_block)
				if err != nil {
					return StyleSheet{}, errors.New("Error while parse block " + string(curr_block) + ": " + err.Error())
				}
				log.Println(block)
				blocks = append(blocks, block)
			}

			curr_block = []rune{}
			continue
		}
	}

	return StyleSheet{Selector: "", AtsInherited: atsInherited, Blocks: blocks, Imports: imports}, nil
}

func parseAt(runes []rune) (at string, err error) {
	fmt.Println("At: ", string(runes))
	return string(at), nil
}
func parseImport(runes []rune) (imp string, err error) {
	tmp_imp := string(runes)
	space_idx := strings.Index(tmp_imp, " ")
	if space_idx == -1 {
		return "", errors.New("Can not parse import: no space symbol")
	}
	tmp_url := tmp_imp[space_idx+1:]
	if len(tmp_url) == 0 {
		return "", errors.New("Can not parse import: no url")
	}

	if tmp_url[len(tmp_url)-1] == ';' {
		tmp_url = tmp_url[:len(tmp_url)-1]
	}

	tmp_url = strings.Trim(tmp_url, "'\"")

	fmt.Println("Import: ", tmp_url)

	return tmp_url, nil
}

func parseBlock(runes []rune) (block Block, err error) {
	rules := make([]Rule, 0)

	tmp_block := string(runes)

	br_idx := strings.Index(tmp_block, "{")

	if br_idx == -1 {
		return Block{}, errors.New("Can not find left bracket on css block")
	}

	sel := strings.Trim(tmp_block[:br_idx], " ")

	lbr_idx := strings.LastIndex(tmp_block, "}")

	if lbr_idx == -1 {
		return Block{}, errors.New("Can not find right bracket on css block")
	}

	cont := strings.Trim(tmp_block[br_idx+1:lbr_idx], " ;")

	var inside_string rune = 0
	buff := []rune{}
	for i, r := range []rune(cont) {
		if inside_string == 0 && r == ';' {
			rule, err := parseRule(buff)
			if err != nil {
				return Block{}, errors.New("Rule error: " + err.Error())
			}
			rules = append(rules, rule)

			buff = []rune{}
			continue
		}
		if inside_string == 0 && (r == '"' || r == '\'') {
			inside_string = r
		} else if inside_string == r {
			inside_string = 0
		}
		buff = append(buff, r)

		if i == len(cont)-1 {
			rule, err := parseRule(buff)
			if err != nil {
				return Block{}, errors.New("Rule error: " + err.Error())
			}
			rules = append(rules, rule)
		}
	}

	return Block{Selector: sel, Rules: rules}, nil
}

func parseAtInherited(runes []rune) (atInherited StyleSheet, err error) {
	fmt.Println("AtInherited: ", string(runes))

	tmp_block := string(runes)

	br_idx := strings.Index(tmp_block, "{")

	if br_idx == -1 {
		return StyleSheet{}, errors.New("Can not find left bracket on css block")
	}

	sel := strings.Trim(tmp_block[:br_idx], " ")

	lbr_idx := strings.LastIndex(tmp_block, "}")

	if lbr_idx == -1 {
		return StyleSheet{}, errors.New("Can not find right bracket on css block")
	}

	cont := strings.Trim(tmp_block[br_idx+1:lbr_idx], " ;")

	fmt.Println("AT SEL: ",sel,", conf: ",cont)
	atInherited,err=Parse(cont)



	fmt.Println("AT: ", atInherited)
	atInherited.Selector = sel
	return

}

func parseRule(runes []rune) (rule Rule, err error) {
	tmp_rule := string(runes)
	sep_idx := strings.Index(tmp_rule, ":")
	if sep_idx == -1 {
		return Rule{}, errors.New("Can not parse rule " + tmp_rule + " can not find ':' symbol ")
	}

	sel := strings.Trim(tmp_rule[0:sep_idx], " ")
	val := strings.Trim(tmp_rule[sep_idx+1:], " ")

	fmt.Println("rule parts: " + sel + "<->" + val + "<->")
	return Rule{Name: sel, Value: val}, nil
}