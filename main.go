package main

import (
	"errors"
	"os"
	"regexp"
	"strconv"
)

// lexer
const (
	WHITESPACE = "\\s"
	NUMBERS    = "\\d"
	PLUS       = "+"
	MINUS      = "-"
	STAR       = "*"
	SLASH      = "/"
	PAR_OPEN   = "("
	PAR_CLOSE  = ")"
)

// lexer
type token struct {
	typ   string
	value string
}

// ast
type astNode struct {
	syntaxType string
	value      string
	left       *astNode
	right      *astNode
}

// ast
var singletonAST *abstractSyntaxTree

type abstractSyntaxTree struct {
	root *astNode
}

var lastOp string      // last operand for ast
var lastSynType string // last syntax type for ast

func main() {
	input := "2*(3+4*2+6)*5*7+2*3"

	tokens, err := lexer(input)
	errCheck(err)

	parser(tokens)
	printAST()

	result := evaluator(getAST().root)
	println("[*] Result:", result)
}

func errCheck(err error) {
	if err != nil {
		println("[!] ERROR:", err.Error())
		os.Exit(1)
	}
}

func lexer(input string) ([]token, error) {
	var tokens []token
	var char string
	var isNumber bool
	var isWhiteSpace bool

	for current := 0; current < len(input); {
		char = string(input[current])                          // take the current char
		isNumber, _ = regexp.MatchString(NUMBERS, char)        // check the current char if is number
		isWhiteSpace, _ = regexp.MatchString(WHITESPACE, char) // check the current char if is whitespace

		if isNumber {
			value := ""

			for isNumber {
				value += char
				current++

				if current >= len(input) {
					break
				}

				char = string(input[current])                   // take the next char
				isNumber, _ = regexp.MatchString(NUMBERS, char) // check the current char if is number
			}

			// if the word is end, append the value as a name
			tokens = append(tokens, token{
				typ:   "number",
				value: value,
			})
			continue
		} else if isWhiteSpace {
			current++
			continue
		} else if char == PAR_OPEN {
			tokens = append(tokens, token{
				typ:   "parOpen",
				value: char,
			})

			current++
			continue
		} else if char == PAR_CLOSE {
			tokens = append(tokens, token{
				typ:   "parClose",
				value: char,
			})

			current++
			continue
		} else if char == PLUS {
			tokens = append(tokens, token{
				typ:   "operation",
				value: char,
			})

			current++
			continue
		} else if char == MINUS {
			tokens = append(tokens, token{
				typ:   "operation",
				value: char,
			})

			current++
			continue
		} else if char == STAR {
			tokens = append(tokens, token{
				typ:   "operation",
				value: char,
			})

			current++
			continue
		} else if char == SLASH {
			tokens = append(tokens, token{
				typ:   "operation",
				value: char,
			})

			current++
			continue
		}

		return nil, errors.New("The \"" + char + "\" is not defined in Lexer")
	}

	return tokens, nil
}

func printTokens(tokens []token) {
	for j := 0; j < len(tokens); j++ {
		println("Token type: ", tokens[j].typ)
		println("Token val: ", tokens[j].value)
	}
}

func getAST() *abstractSyntaxTree {
	if singletonAST == nil {
		singletonAST = &abstractSyntaxTree{
			root: nil,
		}
	}

	return singletonAST
}

func printNextNode(node *astNode, indent int) {
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
	print("> Node Syntax Type: " + node.syntaxType + "\tNode Value: " + node.value + "\n")

	printNextNode(node.left, indent+4)
	printNextNode(node.right, indent+4)
}

func printAST() {
	printNextNode(getAST().root, 0)
}

// Consume target token
func consumeToken(tokens []token, index int) []token {
	oLen := len(tokens)
	newT := tokens[0:index]
	newT = append(newT, tokens[index+1:oLen]...)
	return newT
}

// Get last binary expression node from Right
func getLastBENodeR(node *astNode) *astNode {
	if lastSynType == "ParBinaryExpression" {
		for node.right != nil && (node.right.syntaxType == "BinaryExpression" || node.right.syntaxType == "ParBinaryExpression") {
			node = node.right
		}
	} else {
		for node.right != nil && (node.right.syntaxType == "BinaryExpression") {
			node = node.right
		}
	}

	return node
}

// Get last binary expression node from Left
func getLastBENodeL(node *astNode) *astNode {
	if lastSynType == "ParBinaryExpression" {
		for node.left != nil && (node.left.syntaxType == "BinaryExpression" || node.left.syntaxType == "ParBinaryExpression") {
			node = node.left
		}
	} else {
		for node.left != nil && (node.left.syntaxType == "BinaryExpression") {
			node = node.left
		}
	}

	return node
}

func parser(tokens []token) {
	var isBeforeNumber bool = false

	for i := 0; i < len(tokens); i++ {
		// Arka arkaya iki adet sayi olmasini engeller (1 2 + 3) gibi
		if tokens[i].typ == "number" && isBeforeNumber == true {
			errCheck(errors.New("Unexpected Type Error! Expected: 'Expression'. Current: 'Number'"))
		}

		if tokens[i].typ == "number" && isBeforeNumber == false {
			isBeforeNumber = true
			continue
		}

		// Parantez icindeki islemler icin. Syntax tipi olarak "PreBinaryExpression" kullanilacak
		if tokens[i].typ == "parOpen" {
			tokens = parseParanthesisExp(tokens, &i, &isBeforeNumber)
		}

		if isBeforeNumber && tokens[i].typ == "operation" && tokens[i+1].typ == "number" {
			parseBinaryExpression(&tokens[i-1], &tokens[i], &tokens[i+1], "BinaryExpression")
			isBeforeNumber = false
		} else if isBeforeNumber && tokens[i].typ == "operation" && tokens[i+1].typ == "parOpen" {
			// 2 * (1+2) gibi islemler icin su sekilde davranmasini sagliycaz -> (1+2) * 2 - Bunun icin de islemin solundaki sayiyi saga vercez
			lastIndex := i                                                                                            // En son biraktigimiz index. Cunku bu index parseParanthesis'ten sonra degisecek
			i = i + 1                                                                                                 // index'i parantezin oldugu index'e kaydirdik
			isBeforeNumber = false                                                                                    //onceki tip'in sayi olmadigini bildirdik
			tokens = parseParanthesisExp(tokens, &i, &isBeforeNumber)                                                 // once parantez icini parse ediyoruz
			parseBinaryExpression(&tokens[lastIndex-1], &tokens[lastIndex], &tokens[lastIndex-1], "BinaryExpression") // Sonra en son kaldigimiz index'teki islemi parse ediyoruz
			isBeforeNumber = true
			i = i - 1 // Olmali!
		} else {
			errCheck(errors.New("Undefined Type Error! Expected: 'Number or Expression' Current: '" + tokens[i].typ + "'"))
		}
	}
}

// [!] It changes tokens and returns new tokens. Do not forget to get new tokens.
func parseParanthesisExp(tokens []token, index *int, isBfrNum *bool) []token {
	i := *index
	isBeforeNumber := *isBfrNum
	tokens = consumeToken(tokens, i)

	for ; tokens[i].typ != "parClose"; i++ {
		if tokens[i].typ == "parOpen" {
			tokens = parseParanthesisExp(tokens, &i, &isBeforeNumber)
		}

		if tokens[i].typ == "number" && isBeforeNumber == false {
			isBeforeNumber = true
			continue
		}

		if isBeforeNumber && tokens[i].typ == "operation" && tokens[i+1].typ == "number" {
			parseBinaryExpression(&tokens[i-1], &tokens[i], &tokens[i+1], "ParBinaryExpression")
			isBeforeNumber = false
		} else {
			errCheck(errors.New("Undefined Type Error! Expected: 'Number or Paranthesis Expression' Current: '" + tokens[i].typ + "'"))
		}
	}

	if tokens[i].typ == "parClose" {
		tokens = consumeToken(tokens, i)
	}

	*index = i
	*isBfrNum = isBeforeNumber
	return tokens
}

func parseBinaryExpression(left, current, right *token, synType string) {
	leftNode, err := parseNumberLiteral(left)
	errCheck(err)

	rightNode, err := parseNumberLiteral(right)
	errCheck(err)

	newNode := astNode{
		syntaxType: synType, // "BinaryExpression" || "ParBinaryExpression"
		value:      current.value,
		left:       leftNode,
		right:      rightNode,
	}

	if getAST().root == nil {
		getAST().root = &newNode
		lastOp = current.value
		lastSynType = newNode.syntaxType

		return
	}

	// Eger son eklenmis exp parantez iciyse ve yeni gelen disindaysa
	if lastSynType == "ParBinaryExpression" && newNode.syntaxType == "BinaryExpression" {
		temp := *getAST().root
		newNode.right = &temp
		newNode.left = rightNode

		getAST().root = &newNode

		lastOp = current.value
		lastSynType = newNode.syntaxType
		return
	}

	lastBE := getLastBENodeR(getAST().root)
	// is Last BE Node got from left side of ast
	isDirectionL := false
	// setup for where new node put
	if lastBE.value == STAR || lastBE.value == SLASH {
		lastBE = getLastBENodeL(getAST().root)
		isDirectionL = true
		if lastBE.value == STAR || lastBE.value == SLASH {
			lastBE = getAST().root
		}
	}

	if (lastOp == STAR || lastOp == SLASH) && (newNode.value == PLUS || newNode.value == MINUS) {
		temp := *lastBE
		newNode.right = &temp
		newNode.left = rightNode

		if lastBE == getAST().root {
			getAST().root = &newNode
		} else {
			lastBE = &newNode
		}

		lastOp = current.value
		lastSynType = newNode.syntaxType
		return
	}

	if (lastOp == PLUS || lastOp == MINUS) && (newNode.value == STAR || newNode.value == SLASH) {
		if isDirectionL {
			temp := lastBE.left
			newNode.right = rightNode
			newNode.left = temp
			lastBE.left = &newNode
		} else {
			temp := lastBE.right
			newNode.left = temp
			newNode.right = rightNode
			lastBE.right = &newNode
		}

		lastOp = current.value
		lastSynType = newNode.syntaxType
		return
	}

	// Eger son eklenmis node ve yeni node'un isaretleri ayni ise (+,- ya da *,/ ikilileri) yeni geleni direkt saga ekliycez
	temp := lastBE.right
	newNode.left = temp // It is same with leftNode
	newNode.right = rightNode
	lastBE.right = &newNode

	lastOp = current.value
	lastSynType = newNode.syntaxType
}

func parseNumberLiteral(token *token) (*astNode, error) {
	if token.typ == "number" {
		return &astNode{
			syntaxType: "NumberLiteral",
			value:      token.value,
			left:       nil,
			right:      nil,
		}, nil
	}

	return nil, errors.New("Unexpected Type Error! Expected: 'NumberLiteral' Current:'" + token.typ + "'")
}

func evaluator(node *astNode) int {
	var right int
	var left int

	if node.syntaxType == "NumberLiteral" {
		val, err := strconv.Atoi(node.value)
		errCheck(err)

		return val
	}

	if node.syntaxType == "BinaryExpression" || node.syntaxType == "ParBinaryExpression" {
		if node.right != nil {
			right = evaluator(node.right)
		}

		if node.left != nil {
			left = evaluator(node.left)
		}

		if node.value == STAR {
			return left * right
		} else if node.value == SLASH {
			return left / right
		} else if node.value == PLUS {
			return left + right
		} else if node.value == MINUS {
			return left - right
		}
	}

	return 0
}
