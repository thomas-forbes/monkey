# Monkey

A programming language built following the books [Writing An Interpreter In Go](https://interpreterbook.com/) and it's sequel [Writing A Compiler In Go](https://compilerbook.com/). The language is a simple and dynamically typed language with a syntax similar to JavaScript.

## Example

```monkey
let fib = fn(n) {
  if (n < 2) {
    n;
  } else {
    fib(n - 1) + fib(n - 2);
  }
};

puts(fib(10)); // 55
```

## Usage

Running:
```bash
# Run the REPL
go run main.go
# Run a Monkey source file
go run main.go path/to/source.monkey
```

Testing:
```bash
# All tests
go test ./...
# Specific test
go test ./runner -run TestSpec/errors
```

## Syntax

### Variables

```monkey
let x = 5; // immutable by default variable
let mut y = 10; // mutable variable
y = y + x;
```

### Operations

```monkey
>> 5 + 10
15
>> 5 - 10
-5
>> 5 * 10
50
>> 10 / 5
2
>> "Hello, " + "world!"
"Hello, world!"
```

### Standard Library

built in functions:

```monkey
>> puts("Hello, world!"); // prints to the console
Hello, world!
null
>> len("Hello");
5
>> append([1, 2], 3); 
[1, 2, 3] // new array
```

### Functions

Functions are first class citizens in Monkey.

```monkey
let add = fn(a, b) {
  a + b;
};

add(5, 10); // 15

// closure support
let getMultiplier = fn(x) {
  fn(y) {
    x * y;
  };
}
let doubler = getMultiplier(2);
doubler(5); // 10
```

### Conditionals

```monkey
if (x > 10) {
  puts("x is greater than 10");
} else if (x == 10) {
  puts("x is equal to 10");
} else {
  puts("x is less than 10");
}
```

### Loops

```monkey
for index, item in ["first", "second", "third"] {
  puts(index, item);
}

for key, value in {"a": 1, "b": 2} {
  puts(key + ": " + value);
}

for i in 0..5 {
  puts(i); // prints 0, 1, 2, 3, 4
}

for condition {
  puts("This will run until the condition is false");
}
```
