package object

import "fmt"

var Builtins = []struct {
	Name    string
	Builtin *Builtin
}{
	{
		"len",
		&Builtin{Fn: func(args ...Object) Object {
			if len(args) != 1 {
				return NewWrongArgumentCount(nil, 1, len(args))
			}

			switch arg := args[0].(type) {
			case *Array:
				return &Integer{Value: int64(len(arg.Elements))}
			case *String:
				return &Integer{Value: int64(len(arg.Value))}
			default:
				return NewBuiltinArgumentType(nil, "len", "ARRAY or STRING", string(args[0].Type()))
			}
		},
		},
	},
	{
		"puts",
		&Builtin{Fn: func(args ...Object) Object {
			for _, arg := range args {
				fmt.Print(arg.Inspect(), " ")
			}
			fmt.Println()
			return NULL
		},
		},
	},
	{
		"append",
		&Builtin{Fn: func(args ...Object) Object {
			if len(args) != 2 {
				return NewWrongArgumentCount(nil, 2, len(args))
			}
			if args[0].Type() != ARRAY_OBJ {
				return NewBuiltinArgumentType(nil, "append", "ARRAY", string(args[0].Type()))
			}
			arr := args[0].(*Array)
			length := len(arr.Elements)
			newElements := make([]Object, length+1)
			copy(newElements, arr.Elements)
			newElements[length] = args[1]
			return &Array{Elements: newElements}
		},
		},
	},
}

func GetBuiltinByName(name string) *Builtin {
	for _, def := range Builtins {
		if def.Name == name {
			return def.Builtin
		}
	}
	return nil
}
