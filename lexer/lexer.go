package lexer

import (
	"errors"
	"go-comp-arithmetic/constants"
	"regexp"
)

// lexer
type Token struct {
	Typ   string
	Value string
}

func Lexer(input string) ([]Token, error) {
	var tokens []Token
	var char string
	var isNumber bool
	var isWhiteSpace bool

	for current := 0; current < len(input); {
		char = string(input[current])                                    // take the current char
		isNumber, _ = regexp.MatchString(constants.NUMBERS, char)        // check the current char if is number
		isWhiteSpace, _ = regexp.MatchString(constants.WHITESPACE, char) // check the current char if is whitespace

		if isNumber {
			value := ""

			for isNumber {
				value += char
				current++

				if current >= len(input) {
					break
				}

				char = string(input[current])                             // take the next char
				isNumber, _ = regexp.MatchString(constants.NUMBERS, char) // check the current char if is number
			}

			// if the word is end, append the value as a name
			tokens = append(tokens, Token{
				Typ:   "number",
				Value: value,
			})
			continue
		} else if isWhiteSpace {
			current++
			continue
		} else if char == constants.PAR_OPEN {
			tokens = append(tokens, Token{
				Typ:   "parOpen",
				Value: char,
			})

			current++
			continue
		} else if char == constants.PAR_CLOSE {
			tokens = append(tokens, Token{
				Typ:   "parClose",
				Value: char,
			})

			current++
			continue
		} else if char == constants.PLUS {
			tokens = append(tokens, Token{
				Typ:   "operation",
				Value: char,
			})

			current++
			continue
		} else if char == constants.MINUS {
			tokens = append(tokens, Token{
				Typ:   "operation",
				Value: char,
			})

			current++
			continue
		} else if char == constants.STAR {
			tokens = append(tokens, Token{
				Typ:   "operation",
				Value: char,
			})

			current++
			continue
		} else if char == constants.SLASH {
			tokens = append(tokens, Token{
				Typ:   "operation",
				Value: char,
			})

			current++
			continue
		}

		return nil, errors.New("The \"" + char + "\" is not defined in Lexer")
	}

	return tokens, nil
}

func PrintTokens(tokens []Token) {
	for j := 0; j < len(tokens); j++ {
		println("Token type: ", tokens[j].Typ)
		println("Token val: ", tokens[j].Value)
	}
}

// Consume target token
func ConsumeToken(tokens []Token, index int) []Token {
	oLen := len(tokens)
	newT := tokens[0:index]
	newT = append(newT, tokens[index+1:oLen]...)
	return newT
}
