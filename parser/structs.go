package parser

import "strings"

//border: 1px solid;
type Rule struct {
	Name  string
	Value string
}

//.someclass { rules }
type RuleSet struct {
	Selector string
	Rules    []Rule
}

func (rs RuleSet) getType() CSSType {
	return CSSType_Ruleset
}

func (rs RuleSet) getChilds() []CSSElement {
	return []CSSElement{}
}

func (rs RuleSet) getValue() interface{} {
	return rs.Rules
}

func (rs RuleSet) String() (res string) {
	res = "Ruleset (" + rs.Selector + ")\n"
	for _, r := range rs.Rules {
		res += "\tRule: " + r.Name + "->" + r.Value + "\n"
	}
	return res
}

func (rs RuleSet) StringCSS() (res string) {
	res = rs.Selector + "{\n"
	for _, r := range rs.Rules {
		res += "\t" + r.Name + ": " + r.Value + ";\n"
	}
	res += "}\n"
	return res
}

type Import string

func (i Import) getType() CSSType {
	return CSSType_Import
}

func (i Import) getChilds() []CSSElement {
	return []CSSElement{}
}

func (i Import) getValue() interface{} {
	return string(i)
}

func (i Import) String() (res string) {
	return "Import: " + string(i)
}

func (i Import) StringCSS() (res string) {
	return "@import \"" + string(i) + "\";\n"
}

//Это - либо корневой элемент (Selector=="")
//Либо At-правило, предполагающее вложенность
type CSSStruct struct {
	Selector string
	Childs   []CSSElement
}

func (ai CSSStruct) getType() CSSType {
	return CSSType_AtInherited
}

func (ai CSSStruct) getChilds() []CSSElement {
	return ai.Childs
}

func (ai CSSStruct) getValue() interface{} {
	return ai.Selector
}

func (ai CSSStruct) String() (res string) {
	if ai.Selector == "" {
		res = "Root: \n"
	} else {
		res = "CSSStruct: (" + ai.Selector + ")\n"
	}
	for _, r := range ai.getChilds() {
		res += "-"
		for _, s := range strings.Split(r.String(), "\n") {
			res += "\t" + s + "\n"
		}
	}
	return res
}

func (ai CSSStruct) StringCSS() (res string) {
	need_tab := ""
	if ai.Selector != "" {
		need_tab = "\t"
		res = ai.Selector + " {\n"
	}
	for _, r := range ai.getChilds() {
		for _, s := range strings.Split(r.StringCSS(), "\n") {
			res += need_tab + s + "\n"
		}
	}
	if ai.Selector != "" {
		res += "}\n"
	}
	return res
}

type AtRule string

func (ar AtRule) getType() CSSType {
	return CSSType_AtRule
}

func (ar AtRule) getChilds() []CSSElement {
	return nil
}

func (ar AtRule) getValue() interface{} {
	return string(ar)
}

func (ar AtRule) String() (res string) {
	return "At: " + string(ar)
}

func (ar AtRule) StringCSS() (res string) {
	return string(ar) + ";\n"
}

//Интерфейс для сохранения порядку следования элементов
type CSSElement interface {
	//Получение типа
	getType() CSSType
	//Получение потомков
	getChilds() []CSSElement
	//Получение значения - для элементов, у которых оно есть
	getValue() interface{}
	//Строковый вывод значения
	String() string
	//Вывод в CSS-формате
	StringCSS() string
}

//Типы элементов, которые могут встретиться
type CSSType int

const (
	//Типа .selector { rules }
	CSSType_Ruleset = iota
	//Всякие at-правила, где внутри может быть еще целый дивный мир рулсетов и других at-правил
	CSSType_AtInherited
	//At-правило импорта
	CSSType_Import
	//Прочие at-правила, которые особой ценности не представляют
	CSSType_AtRule
)
