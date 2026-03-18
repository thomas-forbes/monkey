package runner

import (
	"bufio"
	"fmt"
	"io"
	"monkey/compiler"
	"monkey/evaluator"
	"monkey/lexer"
	"monkey/object"
	"monkey/parser"
	"monkey/vm"
	"os"
)

type Engine string

const (
	INTERPRETER Engine = "interpreter"
	VM          Engine = "vm"
)

const PROMPT = ">> "

const ENGINE = VM

func StartRepl(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)
	env := loadStdLib()
	fmt.Printf("Running Monkey REPL with %s engine\n", ENGINE)
	for {
		fmt.Fprintf(out, PROMPT)
		scanned := scanner.Scan()
		if !scanned {
			return
		}
		line := scanner.Text()
		result, _, err := evalProgram(line, env)
		if err != nil {
			fmt.Printf("Error: %s\n", err)
			continue
		}
		if result != nil {
			io.WriteString(out, result.Inspect())
			io.WriteString(out, "\n")
		}
	}
}

func StartFile(fileName string) {
	env := loadStdLib()
	data, err := os.ReadFile(fileName)
	if err != nil {
		fmt.Printf("Couldn't read file: %s\n", err)
	}

	program := string(data)
	_, _, err = evalProgram(program, env)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
	}
}

func evalProgram(program_string string, env *object.Environment) (object.Object, *object.Environment, error) {
	l := lexer.New(program_string)
	p := parser.New(l)
	program := p.ParseProgram()
	if len(p.Errors()) != 0 {
		printParserErrors(os.Stdout, p.Errors())
		return nil, nil, fmt.Errorf("Parser errors encountered")
	}

	var result object.Object
	switch ENGINE {
	case INTERPRETER:
		result = evaluator.Eval(program, env)
		result = evaluator.Eval(program, env)
	case VM:
		comp := compiler.New()
		err := comp.Compile(program)
		if err != nil {
			return nil, nil, fmt.Errorf("Woops! Compilation failed:\n %s\n", err)
		}
		machine := vm.New(comp.Bytecode())
		err = machine.Run()

		if err != nil {
			return nil, nil, fmt.Errorf("Woops! Executing bytecode failed:\n %s\n", err)
		}
		result = machine.StackTop()
	default:
		return nil, nil, fmt.Errorf("Unknown engine: %s", ENGINE)
	}

	if evaluator.IsError(result) {
		return nil, nil, fmt.Errorf("Evaluation error: %s", result.Inspect())
	}

	return result, env, nil
}

func printParserErrors(out io.Writer, errors []string) {
	for _, msg := range errors {
		io.WriteString(out, "\t"+msg+"\n")
	}
}

func loadStdLib() *object.Environment {
	_, env, err := evalProgram(stdlib, object.NewEnvironment())
	if err != nil {
		panic(fmt.Errorf("Error loading stdlib: %s\n", err))
	}
	return env
}
