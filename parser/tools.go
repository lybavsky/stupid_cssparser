package parser

//Поиск по названию атрибута
func (ss *StyleSheet) FindByKey(key string, callback func(rule *Rule)) {
	cstr := ss.Model
	rules := getByKey(cstr, key)

	for _, r := range rules {
		callback(r)
	}

}

func getByKey(cssStruct CSSStruct, key string) (rules []*Rule) {
	for _, el := range *cssStruct.Childs {
		switch el.getType() {
		case CSSType_Ruleset:
			{
				var ruls = el.getValue().([]*Rule)
				for _, r := range ruls {
					if r.Name == key {
						rules = append(rules, r)
					}
				}
				break
			}
		case CSSType_AtInherited:
			{
				rules = append(rules, getByKey(el.(CSSStruct), key)...)
			}
		}
	}
	return rules
}

func (rule Rule) Delete() {
	par := rule.Parent
	//rs := (*par).Rules
	//found := false
	//i := -1
	(*par).Mux.Lock()

	par.Rules = []*Rule{}
	//for i = 0; i < len(rs); i++ {
	//	if *rs[i] == rule {
	//		found = true
	//		fmt.Println("FOUND TO DELETE idx ", i)
	//		break
	//	}
	//}
	//if found {
	//	copy(rs[i:], rs[i+1:]) // Shift a[i+1:] left one index.
	//	rs = rs[:len(rs)-1]    // Truncate slice.
	//	(*par).Rules = []*Rule{}
	//}
	(*par).Mux.Unlock()

}
