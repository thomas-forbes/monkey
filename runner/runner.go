package runner

import (
	"bufio"
	"fmt"
	"io"
	"monkey/evaluator"
	"monkey/lexer"
	"monkey/object"
	"monkey/parser"
	"os"
)

const PROMPT = ">> "

func StartRepl(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)
	env := loadStdLib()
	for {
		fmt.Fprintf(out, PROMPT)
		scanned := scanner.Scan()
		if !scanned {
			return
		}
		line := scanner.Text()
		l := lexer.New(line)
		p := parser.New(l)
		program := p.ParseProgram()
		if len(p.Errors()) != 0 {
			printParserErrors(out, p.Errors())
			continue
		}
		evaluated := evaluator.Eval(program, env)
		if evaluated != nil {
			io.WriteString(out, evaluated.Inspect())
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
	evalProgram(program, env)
}

func evalProgram(program_string string, env *object.Environment) (*object.Environment, error) {
	l := lexer.New(program_string)
	p := parser.New(l)
	program := p.ParseProgram()
	if len(p.Errors()) != 0 {
		printParserErrors(os.Stdout, p.Errors())
		return nil, fmt.Errorf("Parser errors encountered")
	}

	evaluated := evaluator.Eval(program, env)
	if evaluator.IsError(evaluated) {
		fmt.Printf("Error: %s\n", evaluated.Inspect())
		return nil, fmt.Errorf("Evaluation error: %s", evaluated.Inspect())
	}
	return env, nil
}

func printParserErrors(out io.Writer, errors []string) {
	for _, msg := range errors {
		io.WriteString(out, "\t"+msg+"\n")
	}
}

func loadStdLib() *object.Environment {
	env, err := evalProgram(stdlib, object.NewEnvironment())
	if err != nil {
		panic(fmt.Errorf("Error loading stdlib: %s\n", err))
	}
	return env
}
