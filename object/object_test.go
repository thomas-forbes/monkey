package object

import (
	"monkey/token"
	"testing"
)

func TestStringHashKey(t *testing.T) {
	hello1 := &String{Value: "Hello World"}
	hello2 := &String{Value: "Hello World"}
	diff1 := &String{Value: "My name is johnny"}
	diff2 := &String{Value: "My name is johnny"}
	if hello1.HashKey() != hello2.HashKey() {
		t.Errorf("strings with same content have different hash keys")
	}
	if diff1.HashKey() != diff2.HashKey() {
		t.Errorf("strings with same content have different hash keys")
	}
	if hello1.HashKey() == diff1.HashKey() {
		t.Errorf("strings with different content have same hash keys")
	}
}

func TestStructuredErrorWrapsDetailAndToken(t *testing.T) {
	tok := token.Token{Type: token.IDENT, Literal: "foobar", Line: 2, Column: 3, Offset: 10}
	errObj := NewUnknownIdentifier(&tok, "foobar")

	if errObj.Error() != "identifier not found: foobar" {
		t.Fatalf("wrong error string. got=%q", errObj.Error())
	}
	if errObj.Inspect() != "ERROR: identifier not found: foobar" {
		t.Fatalf("wrong inspect string. got=%q", errObj.Inspect())
	}
	detail, ok := errObj.Detail.(UnknownIdentifier)
	if !ok {
		t.Fatalf("wrong detail type. got=%T", errObj.Detail)
	}
	if detail.Name != "foobar" {
		t.Fatalf("wrong detail payload. got=%q", detail.Name)
	}
	if errObj.Tok == nil || errObj.Tok.Line != 2 || errObj.Tok.Column != 3 {
		t.Fatalf("wrong token stored. got=%+v", errObj.Tok)
	}
}
