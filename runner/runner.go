package runner

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"monkey/compiler"
	"monkey/evaluator"
	"monkey/lexer"
	"monkey/object"
	"monkey/parser"
	"monkey/vm"
	"os"
	"time"
)

type Engine string

const (
	INTERPRETER Engine = "eval"
	VM          Engine = "vm"
)

const PROMPT = ">> "

var engineFlag = flag.String("engine", "vm", "use 'vm' or 'eval'")

func GetEngine() (Engine, error) {
	switch *engineFlag {
	case "vm":
		return VM, nil
	case "eval":
		return INTERPRETER, nil
	default:
		return "", fmt.Errorf("Unknown engine: %s", *engineFlag)
	}
}

func StartRepl(engine Engine, in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)
	env := loadStdLib(engine)
	fmt.Printf("Running Monkey REPL with %s engine\n", engine)
	for {
		fmt.Fprintf(out, PROMPT)
		scanned := scanner.Scan()
		if !scanned {
			return
		}
		line := scanner.Text()
		result, _, _, err := evalProgram(engine, line, env)
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

func StartFile(engine Engine, fileName string) {
	env := loadStdLib(engine)
	data, err := os.ReadFile(fileName)
	if err != nil {
		fmt.Printf("Couldn't read file: %s\n", err)
	}

	program := string(data)
	result, _, duration, err := evalProgram(engine, program, env)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
	}
	fmt.Printf("engine=%s, result=%s, duration=%s\n", engine, result.Inspect(), duration)
}

func evalProgram(engine Engine, program_string string, env *object.Environment) (object.Object, *object.Environment, *time.Duration, error) {
	l := lexer.New(program_string)
	p := parser.New(l)
	program := p.ParseProgram()
	if len(p.Errors()) != 0 {
		printParserErrors(os.Stdout, p.Errors())
		return nil, nil, nil, fmt.Errorf("Parser errors encountered")
	}

	var result object.Object
	var duration time.Duration

	switch engine {
	case INTERPRETER:
		start := time.Now()
		result = evaluator.Eval(program, env)
		duration = time.Since(start)
	case VM:
		comp := compiler.New()
		err := comp.Compile(program)
		if err != nil {
			return nil, nil, nil, fmt.Errorf("Woops! Compilation failed:\n %s\n", err)
		}
		machine := vm.New(comp.Bytecode())

		start := time.Now()
		err = machine.Run()

		if err != nil {
			return nil, nil, nil, fmt.Errorf("Woops! Executing bytecode failed:\n %s\n", err)
		}

		duration = time.Since(start)

		result = machine.LastPoppedStackElem()
	default:
		return nil, nil, nil, fmt.Errorf("Unknown engine: %s", engine)
	}

	if evaluator.IsError(result) {
		return nil, nil, nil, fmt.Errorf("Evaluation error: %s", result.Inspect())
	}

	return result, env, &duration, nil
}

func printParserErrors(out io.Writer, errors []string) {
	for _, msg := range errors {
		io.WriteString(out, "\t"+msg+"\n")
	}
}

func loadStdLib(engine Engine) *object.Environment {
	_, env, _, err := evalProgram(engine, stdlib, object.NewEnvironment())
	if err != nil {
		panic(fmt.Errorf("Error loading stdlib: %s\n", err))
	}
	return env
}
