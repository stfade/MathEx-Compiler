package evaluator

import (
	"go-comp-arithmetic/constants"
	"go-comp-arithmetic/errs"
	"go-comp-arithmetic/parser"
	"strconv"
)

func Evaluate() int {
	return evaluator(parser.GetASTRoot())
}

func evaluator(node *parser.AstNode) int {
	var right int
	var left int

	if node.SyntaxType == "NumberLiteral" {
		val, err := strconv.Atoi(node.Value)
		errs.ErrCheck(err)

		return val
	}

	if node.SyntaxType == "BinaryExpression" || node.SyntaxType == "ParBinaryExpression" {
		if node.Right != nil {
			right = evaluator(node.Right)
		}

		if node.Left != nil {
			left = evaluator(node.Left)
		}

		if node.Value == constants.STAR {
			return left * right
		} else if node.Value == constants.SLASH {
			return left / right
		} else if node.Value == constants.PLUS {
			return left + right
		} else if node.Value == constants.MINUS {
			return left - right
		}
	}

	return 0
}
