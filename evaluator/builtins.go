package evaluator

import (
	"monkey/object"
)

var builtins = map[string]*object.Builtin{
	"len":    object.GetBuiltinByName("len"),
	"append": object.GetBuiltinByName("append"),
	"puts":   object.GetBuiltinByName("puts"),
}
