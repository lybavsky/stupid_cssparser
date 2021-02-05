package parser

//.someclass { rules }
type Block struct {
	Selector string
	Rules    []Rule
}

//border: 1px solid;
type Rule struct {
	Name  string
	Value string
}

//Блок и корневой (когда селектор пуст) и наследуемого at-контента (вида @..)
type StyleSheet struct {
	Selector string

	Blocks []Block
	//@media
	AtsInherited []StyleSheet
	//@import
	Imports []string
	//@
	Ats []string
}
