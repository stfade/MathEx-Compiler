package translator

import (
	"fmt"
	"go-comp-arithmetic/constants"
	"go-comp-arithmetic/errs"
	"go-comp-arithmetic/parser"
	"os"
	"strconv"
)

const filename = "result.txt"

var file *os.File

func openFile() {
	var err error

	_, err = os.Stat(filename)
	if err == nil {
		// If the file exists then remove the file
		err = os.Remove(filename)
		errs.ErrCheck(err)
	}

	file, err = os.OpenFile(
		filename,
		os.O_APPEND|os.O_WRONLY|os.O_CREATE,
		0666)

	errs.ErrCheck(err)
}

func writeToFile(input string) {
	byteSlice := []byte(input + "\n")
	_, err := file.Write(byteSlice)
	errs.ErrCheck(err)
}

func closeFile() {
	file.Close()
}

// It writes the mathematical expressions as words to a text file (result.txt)
func Translate() {
	openFile()
	defer closeFile()

	if parser.GetASTRoot() != nil {
		translator(parser.GetASTRoot())
	}

}

func translator(node *parser.AstNode) int {
	var right int
	var left int

	if node.SyntaxType == "NumberLiteral" {
		val, err := strconv.Atoi(node.Value)
		errs.ErrCheck(err)

		return val
	}

	if node.SyntaxType == "BinaryExpression" || node.SyntaxType == "ParBinaryExpression" {
		if node.Right != nil {
			right = translator(node.Right)
		}

		if node.Left != nil {
			left = translator(node.Left)
		}

		if node.Value == constants.STAR {
			writeToFile("MULTIPLY " + fmt.Sprint(left) + " WITH " + fmt.Sprint(right))
			return left * right

		} else if node.Value == constants.SLASH {
			writeToFile("DIVIDE " + fmt.Sprint(left) + " TO " + fmt.Sprint(right))
			return left / right

		} else if node.Value == constants.PLUS {
			writeToFile("ADD " + fmt.Sprint(left) + " TO " + fmt.Sprint(right))
			return left + right

		} else if node.Value == constants.MINUS {
			writeToFile("SUBTRACT " + fmt.Sprint(right) + " FROM " + fmt.Sprint(left))
			return left - right
		}
	}

	return 0
}
