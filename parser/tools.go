package parser

import (
	"fmt"
	"log"
)

//Поиск по названию атрибута
func (ss *StyleSheet) FindByKey(key string, callback func(rule *Rule)) {
	cstr := &ss.Model
	rules := getByKey(cstr, key)

	log.Println(rules)

	for _, r := range rules {
		if r.Parent != nil {
			callback(r)
		}
	}

}

func getByKey(cssStruct *CSSStruct, key string) (rules []*Rule) {
	for _, el := range cssStruct.Childs {
		switch el.getType() {
		case CSSType_Ruleset:
			{
				var ruls = el.(*RuleSet).Rules
				for _, r := range ruls {
					if r.Name == key || key == "" {
						rules = append(rules, r)
					}
				}
				break
			}
		case CSSType_AtInherited:
			{

				ch := el.(*CSSStruct)
				rules = append(rules, getByKey(ch, key)...)
			}
		}
	}
	return rules
}

func (rule *Rule) Delete() {
	par := rule.Parent
	rs := (*par).Rules
	found := false
	i := -1
	(*par).Mux.Lock()

	par.Rules = []*Rule{}
	for i = 0; i < len(rs); i++ {
		if rs[i] == rule {
			found = true
			fmt.Println("FOUND TO DELETE idx ", i)
			break
		}
	}
	if found {
		copy(rs[i:], rs[i+1:]) // Shift a[i+1:] left one index.
		rs = rs[:len(rs)-1]    // Truncate slice.
		(*par).Rules = rs
	}
	(*par).Mux.Unlock()

	//rule.Parent.Rules = []*Rule{}
	//log.Println(rule.Parent)
}
