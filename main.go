package main

import (
	"fmt"
	"go-comp-arithmetic/errs"
	"go-comp-arithmetic/evaluator"
	"go-comp-arithmetic/lexer"
	"go-comp-arithmetic/parser"
	"go-comp-arithmetic/translator"
)

func main() {
	runTest("(1+2*3+4)*(5+6*7+8*9)")
	// runEvalTestCase("(1+2*3+4)*(5+6*7+8*9)")
	// runAllEvalTestCases()
}

func runEvalTestCase(input string) {
	tokens, err := lexer.Lexer(input)
	errs.ErrCheck(err)

	parser.Parse(tokens)
	parser.PrintAST()

	result := evaluator.Evaluate()
	println("[*] Result:", result)
}

func runAllEvalTestCases() {
	input := "1+2*3+4"
	testEvaluator(input, 11)

	input = "(1+2*3+4)"
	testEvaluator(input, 11)

	input = "(1+2)*(3)"
	testEvaluator(input, 9)

	input = "(3+4*2+6)*2"
	testEvaluator(input, 34)

	input = "3*(1+2)"
	testEvaluator(input, 9)

	input = "2*(3+4*2+6)*5*7+2*3"
	testEvaluator(input, 1196)

	input = "(1+2)*(3+4)"
	testEvaluator(input, 21)

	input = "(1+2)*(3+4)*(5+6)"
	testEvaluator(input, 231)

	input = "5*(1+2)*(3+4)*7+8*9+10+15"
	testEvaluator(input, 832)

	input = "(1+2*3+5)*(3+4*2+3*2)"
	testEvaluator(input, 204)
}

func testEvaluator(input string, expectedRes int) {
	tokens, err := lexer.Lexer(input)
	errs.ErrCheck(err)

	parser.Parse(tokens)
	result := evaluator.Evaluate()

	if result == expectedRes {
		println("\u2713 ::DONE " + input + " = " + fmt.Sprint(result) + " is correct! Expected: " + fmt.Sprint(expectedRes))
	} else {
		println("\u26A0 ::ERROR " + input + " = " + fmt.Sprint(result) + " is not correct! Expected: " + fmt.Sprint(expectedRes))
	}

	parser.DeleteASTRoot()
}

func runTest(input string) {
	testTranslator(input)
	parser.DeleteASTRoot()
}

func testTranslator(input string) {
	tokens, err := lexer.Lexer(input)
	errs.ErrCheck(err)

	parser.Parse(tokens)
	translator.Translate()
}
