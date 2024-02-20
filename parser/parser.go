package parser

import (
	"errors"
	"go-comp-arithmetic/constants"
	"go-comp-arithmetic/errs"
	"go-comp-arithmetic/lexer"
)

// Abstract Syntax Tree Node
type AstNode struct {
	SyntaxType string
	Value      string
	Left       *AstNode
	Right      *AstNode
}

// [!] AST must be singleton
var singletonAST *abstractSyntaxTree

type abstractSyntaxTree struct {
	root *AstNode
}

var lastOp string      // last operand for ast
var lastSynType string // last syntax type for ast

func getAST() *abstractSyntaxTree {
	if singletonAST == nil {
		singletonAST = &abstractSyntaxTree{
			root: nil,
		}
	}

	return singletonAST
}

func GetASTRoot() *AstNode {
	return getAST().root
}

func DeleteASTRoot() {
	getAST().root = nil
}

func printNextNode(node *AstNode, indent int) {
	if node == nil {
		return
	}

	// indent
	print("\n")
	for i := 0; i < indent-2; i++ {
		print(" ")
	}

	print("└")
	for i := 0; i < 2; i++ {
		print("─")
	}

	// Print current first
	print("> Node Syntax Type: " + node.SyntaxType + "\tNode Value: " + node.Value + "\n")

	printNextNode(node.Left, indent+4)
	printNextNode(node.Right, indent+4)
}

func PrintAST() {
	printNextNode(getAST().root, 0)
}

func newAstNode(left, current, right *lexer.Token, synType string) *AstNode {
	leftNode, err := parseNumberLiteral(left)
	errs.ErrCheck(err)

	rightNode, err := parseNumberLiteral(right)
	errs.ErrCheck(err)

	return &AstNode{
		SyntaxType: synType, // "BinaryExpression" || "ParBinaryExpression"
		Value:      current.Value,
		Left:       leftNode,
		Right:      rightNode,
	}
}

func printAstNode(node *AstNode) {
	println("Node.type: " + node.SyntaxType)
	println("Node.val: " + node.Value)
	println("Node.l.val: " + node.Left.Value)
	println("Node.r.val: " + node.Right.Value)
}

// Get last binary expression node from the right side
func getLastBENodeR(node *AstNode) *AstNode {
	if lastSynType == "ParBinaryExpression" {
		for node.Right != nil && (node.Right.SyntaxType == "BinaryExpression" || node.Right.SyntaxType == "ParBinaryExpression") {
			node = node.Right
		}
	} else {
		for node.Right != nil && (node.Right.SyntaxType == "BinaryExpression") {
			node = node.Right
		}
	}

	return node
}

// Get last binary expression node from Left
func getLastBENodeL(node *AstNode) *AstNode {
	if lastSynType == "ParBinaryExpression" {
		for node.Left != nil && (node.Left.SyntaxType == "BinaryExpression" || node.Left.SyntaxType == "ParBinaryExpression") {
			node = node.Left
		}
	} else {
		for node.Left != nil && (node.Left.SyntaxType == "BinaryExpression") {
			node = node.Left
		}
	}

	return node
}

// Get last binary expression node which op is plus or minus from Right
func getLastBEPlusOrMinusR(node *AstNode) *AstNode {
	if lastSynType == "ParBinaryExpression" {
		for node.Right != nil && (node.Right.SyntaxType == "BinaryExpression" || node.Right.SyntaxType == "ParBinaryExpression") && (node.Right.Value == "+" || node.Right.Value == "-") {
			node = node.Right
		}
	} else {
		for node.Right != nil && node.Right.SyntaxType == "BinaryExpression" && (node.Right.Value == "+" || node.Right.Value == "-") {
			node = node.Right
		}
	}

	return node
}

// Get last binary expression node which op is plus or minus from Left
func getLastBEPlusOrMinusL(node *AstNode) *AstNode {
	if lastSynType == "ParBinaryExpression" {
		for node.Left != nil && (node.Left.SyntaxType == "BinaryExpression" || node.Left.SyntaxType == "ParBinaryExpression") && (node.Left.Value == "+" || node.Left.Value == "-") {
			node = node.Left
		}
	} else {
		for node.Left != nil && node.Left.SyntaxType == "BinaryExpression" && (node.Left.Value == "+" || node.Left.Value == "-") {
			node = node.Left
		}
	}

	return node
}

func getParentNode(child *AstNode) *AstNode {
	if child == getAST().root {
		return nil
	}

	node := getAST().root

	for node.Right != nil {
		for node.Left != nil {
			if child == node.Left {
				return node
			}

			node = node.Left
		}

		if child == node.Right {
			return node
		}

		node = node.Right
	}

	for node.Left != nil {
		for node.Right != nil {
			if child == node.Right {
				return node
			}

			node = node.Right
		}

		if child == node.Left {
			return node
		}

		node = node.Left
	}

	return nil
}

func Parse(tokens []lexer.Token) {
	var isBeforeNumber bool = false

	for i := 0; i < len(tokens); i++ {
		if i >= len(tokens) {
			break
		}

		// Arka arkaya iki adet sayi olmasini engeller (1 2 + 3) gibi
		if tokens[i].Typ == "number" && isBeforeNumber == true {
			errs.ErrCheck(errors.New("Unexpected Type Error! Expected: 'Expression'. Current: 'Number'"))
		}

		if tokens[i].Typ == "number" && isBeforeNumber == false {
			isBeforeNumber = true
			continue
		}

		// Parantez icindeki islemler icin. Syntax tipi olarak "PreBinaryExpression" kullanilacak
		if tokens[i].Typ == "parOpen" {
			tokens = parseParanthesisExp(tokens, &i, &isBeforeNumber)
			continue
		}

		if tokens[i-1].Typ == "number" && tokens[i].Typ == "operation" && tokens[i+1].Typ == "number" {
			node := newAstNode(&tokens[i-1], &tokens[i], &tokens[i+1], "BinaryExpression")
			parseBinaryExpression(node)

			isBeforeNumber = false
		} else if !isBeforeNumber && lastSynType == "ParBinaryExpression" && tokens[i].Typ == "operation" && tokens[i+1].Typ == "parOpen" {
			node := newAstNode(&tokens[i-1], &tokens[i], &tokens[i+2], "BinaryExpression")
			parseBinaryExpression(node)

			i = i + 1
			isBeforeNumber = false
			tokens = parseParanthesisExp(tokens, &i, &isBeforeNumber)

		} else if isBeforeNumber && tokens[i].Typ == "operation" && tokens[i+1].Typ == "parOpen" {
			// 2 * (1+2) gibi islemler icin su sekilde davranmasini sagliycaz -> (1+2) * 2 - Bunun icin de islemin solundaki sayiyi saga vercez
			lastIndex := i                                            // En son biraktigimiz index. Cunku bu index parseParanthesis'ten sonra degisecek
			i = i + 1                                                 // index'i parantezin oldugu index'e kaydirdik
			isBeforeNumber = false                                    //onceki tip'in sayi olmadigini bildirdik
			tokens = parseParanthesisExp(tokens, &i, &isBeforeNumber) // once parantez icini parse ediyoruz

			node := newAstNode(&tokens[lastIndex-1], &tokens[lastIndex], &tokens[lastIndex-1], "BinaryExpression")
			parseBinaryExpression(node) // Sonra en son kaldigimiz index'teki islemi parse ediyoruz

			lastSynType = "ParBinaryExpression"
			// isBeforeNumber = false
			// i = i - 1 // Olmali!
		} else {
			errs.ErrCheck(errors.New("Undefined Type Error! Expected: 'Number or Expression' Current: '" + tokens[i].Typ + "'"))
		}
	}
}

// [!] It changes tokens and returns new tokens. Do not forget to get new tokens.
func parseParanthesisExp(tokens []lexer.Token, index *int, isBfrNum *bool) []lexer.Token {
	i := *index
	isBeforeNumber := *isBfrNum
	tokens = lexer.ConsumeToken(tokens, i)

	for ; tokens[i].Typ != "parClose"; i++ {
		if tokens[i].Typ == "parOpen" {
			tokens = parseParanthesisExp(tokens, &i, &isBeforeNumber)
		}

		if tokens[i].Typ == "number" && isBeforeNumber == false {
			isBeforeNumber = true
			continue
		}

		if isBeforeNumber && tokens[i].Typ == "operation" && tokens[i+1].Typ == "number" {
			node := newAstNode(&tokens[i-1], &tokens[i], &tokens[i+1], "ParBinaryExpression")
			parseBinaryExpression(node)
			isBeforeNumber = false
		} else {
			errs.ErrCheck(errors.New("Undefined Type Error! Expected: 'Number or Paranthesis Expression' Current: '" + tokens[i].Typ + "'"))
		}
	}

	if tokens[i].Typ == "parClose" {
		tokens = lexer.ConsumeToken(tokens, i)
		isBeforeNumber = false
	}

	*index = i - 1
	*isBfrNum = isBeforeNumber
	return tokens
}

func parseBinaryExpression(newNode *AstNode) {
	if getAST().root == nil {
		getAST().root = newNode

		lastOp = newNode.Value
		lastSynType = newNode.SyntaxType
		return
	}

	// Eger son eklenmis exp parantez iciyse ve yeni gelen disindaysa
	// If last appended exp is in any paranthesis exp and the new one is outside of the paranthesis exp
	if lastSynType == "ParBinaryExpression" && newNode.SyntaxType == "BinaryExpression" {
		// Root'u degistiriyor
		temp := *getAST().root
		newNode.Left = newNode.Right
		newNode.Right = &temp

		getAST().root = newNode

		lastOp = newNode.Value
		lastSynType = newNode.SyntaxType

		return
	}

	lastBE := getLastBENodeR(getAST().root)
	// is Last BE Node got from left side of ast
	isDirectionL := false
	// setup for where new node put
	if lastBE.Value == constants.STAR || lastBE.Value == constants.SLASH {
		lastBE = getLastBENodeL(getAST().root)
		isDirectionL = true
		if lastBE.Value == constants.STAR || lastBE.Value == constants.SLASH {
			lastBE = getAST().root
		}
	}

	if lastSynType == "BinaryExpression" && newNode.SyntaxType == "ParBinaryExpression" {
		newNode.Left = getAST().root.Left
		getAST().root.Left = newNode
		// newNode.left = lastBE.left
		// lastBE.left = newNode

		lastOp = newNode.Value
		lastSynType = newNode.SyntaxType

		return
	}

	if (lastOp == constants.STAR || lastOp == constants.SLASH) && (newNode.Value == constants.PLUS || newNode.Value == constants.MINUS) {
		if isDirectionL {
			lastBE = getLastBEPlusOrMinusL(getAST().root)
		} else {
			lastBE = getLastBEPlusOrMinusR(getAST().root)
		}

		temp := *lastBE

		if lastBE == getAST().root {
			newNode.Left = newNode.Right
			newNode.Right = &temp
			getAST().root = newNode
		} else {
			newNode.Left = &temp
			parent := getParentNode(lastBE)

			if isDirectionL {
				parent.Left = newNode
			} else {
				parent.Right = newNode
			}
		}

		lastOp = newNode.Value
		lastSynType = newNode.SyntaxType

		return
	}

	if (lastOp == constants.PLUS || lastOp == constants.MINUS) && (newNode.Value == constants.STAR || newNode.Value == constants.SLASH) {
		if isDirectionL {
			temp := lastBE.Left

			if lastBE == getAST().root && (lastBE.Value == "*" || lastBE.Value == "/") {
				temp.Right = newNode
			} else if lastBE == getAST().root && (lastBE.Value == "+" || lastBE.Value == "-") {
				newNode.Left = temp
				lastBE.Left = newNode
			} else {
				lastBE.Right = temp
				lastBE.Left = newNode
			}
		} else {
			temp := *lastBE.Right
			newNode.Left = &temp
			lastBE.Right = newNode
		}

		lastOp = newNode.Value
		lastSynType = newNode.SyntaxType

		return
	}

	// Eger son eklenmis node ve yeni node'un isaretleri ayni ise (+,- ya da *,/ ikilileri) yeni geleni direkt saga ekliycez
	temp := lastBE.Right
	newNode.Left = temp // It is same with leftNode
	lastBE.Right = newNode

	lastOp = newNode.Value
	lastSynType = newNode.SyntaxType
}

func parseNumberLiteral(token *lexer.Token) (*AstNode, error) {
	if token.Typ == "number" {
		return &AstNode{
			SyntaxType: "NumberLiteral",
			Value:      token.Value,
			Left:       nil,
			Right:      nil,
		}, nil
	}

	return nil, errors.New("Unexpected Type Error! Expected: 'NumberLiteral' Current:'" + token.Typ + "'")
}
