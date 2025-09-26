package bsd_testtool

import (
	"encoding/hex"
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"unicode"
)

func FmtSprintf(format string, ctx *ActionContext) string {
	re := regexp.MustCompile(`\{(\d+|\d:\S+)\}`)
	return re.ReplaceAllStringFunc(format, func(s string) string {
		if ctx == nil {
			return "[ctx nil]"
		}
		idx := re.FindStringSubmatch(s)[1]
		arrIndexBegin, arrIndexEnd := -1, -1
		if strings.Contains(idx, ":") {
			arrFmtStrings := strings.Split(idx, ":")
			idx = arrFmtStrings[0]
			arrIndexStrings := strings.Split(arrFmtStrings[1], ",")
			if len(arrIndexStrings) == 2 {
				arrIndexEnd, _ = strconv.Atoi(arrIndexStrings[1])
			}
			arrIndexBegin, _ = strconv.Atoi(arrIndexStrings[0])
		}
		switch idx {
		case "0":
			return ctx.LastModuleName
		case "1":
			if arrIndexBegin > len(ctx.LastSerialBuffer)-1 || arrIndexEnd > len(ctx.LastSerialBuffer)-1 {
				return "[index out of range]"
			}
			if arrIndexBegin != -1 && arrIndexEnd == -1 {
				return hex.EncodeToString(ctx.LastSerialBuffer[arrIndexBegin : arrIndexBegin+1])
			} else if arrIndexBegin != -1 {
				return hex.EncodeToString(ctx.LastSerialBuffer[arrIndexBegin:arrIndexEnd])
			} else {
				return hex.EncodeToString(ctx.LastSerialBuffer)
			}
		case "2":
			return ctx.LastExecResult.Error()
		}
		return "[nil]"
	})
}

func FmtGetVar(format string, ctx *ActionContext) (varType string, v any) {
	return
}

func FmtEvalCondition(format string, ctx *ActionContext) (bool, error) {
	return true, nil
}

var ErrVarNotClosed error = errors.New("var not closed")
var ErrInvalidBinaryOp error = errors.New("invalid binary operator")
var ErrInvalidToken error = errors.New("invalid token")

func TestTokenize(input string) ([]token, error) {
	return tokenize(input)
}

func tokenize(input string) ([]token, error) {
	var tokens []token
	i := 0
	for i < len(input) {
		ch := input[i]

		// 跳过空格
		if unicode.IsSpace(rune(ch)) {
			i++
			continue
		}

		if ch == '(' {
			tokens = append(tokens, token{string(ch), TokenLB})
			i++
			continue
		}

		if ch == ')' {
			tokens = append(tokens, token{string(ch), TokenRB})
			i++
			continue
		}

		// 变量
		// {0} 上下文里的 LastModuleName
		// {1} 上下文里的 LastSerialBuffer
		// {2} 上下文里的 LastExecResult
		// {var} 变量
		// {array: 0} -> array[0]
		// {array: 0, 10} -> array[0:10]
		if ch == '{' {
			j := i + 1
			for j < len(input) && input[j] != '}' {
				j++
			}
			if j < len(input) {
				tokens = append(tokens, token{input[i : j+1], TokenVar})
				i = j + 1
				continue
			} else {
				return nil, ErrVarNotClosed
			}
		}

		// 整数(支持0x开头的十六进制数字)
		if unicode.IsDigit(rune(ch)) {
			j := i + 1

			// 0x
			if ch == '0' && j < len(input) && (input[j] == 'x' || input[j] == 'X') {
				j++
				for j < len(input) && isHexDigit(input[j]) {
					j++
				}
			} else {
				for j < len(input) && unicode.IsDigit(rune(input[j])) {
					j++
				}
			}

			tokens = append(tokens, token{input[i:j], TokenConst})
			i = j
			continue
		}

		// 运算符
		if isOperatorChar(ch) {
			// 尝试识别两字符运算符
			if i+1 < len(input) && isOperatorChar(input[i+1]) {
				op := input[i : i+2]
				// 只接受合法双字符运算符
				if isTwoCharOp(op) {
					tokens = append(tokens, token{op, getOpTokenType(op)})
					i += 2
					continue
				}
			}
			if !isValidSingleOpChar(ch) {
				return nil, ErrInvalidToken
			}
			// 单字符运算符
			tokens = append(tokens, token{string(ch), getOpTokenType(string(ch))})
			i++
			continue
		}

		return nil, ErrInvalidToken

	}

	return tokens, nil
}

type nodeType int

const (
	NodeOp nodeType = iota
	NodeVar
	NodeConst
)

type tokenType int

const (
	TokenOr  tokenType = iota // ||
	TokenAnd                  // &&

	TokenLt    // <
	TokenLEt   // <=
	TokenMt    // >
	TokenMEt   // >=
	TokenEq    // ==
	TokenNotEq // !=

	TokenVar   // {var}
	TokenConst // 0x0A | 10
	TokenLB    // (
	TokenRB    // )
)

func (t tokenType) isCompare() bool {
	return (t >= TokenLt && t <= TokenNotEq)
}

type AstNode struct {
	Type  nodeType
	Value string
	Left  *AstNode
	Right *AstNode
}

type token struct {
	t      string
	toType tokenType
}

type parser struct {
	t   []token
	idx int
}

/*
优先级
低到高
0. ||
1. &&
2. ==, !=, >, <, >=, <=
3. 括号

{test} > 0x10 || {1:5} == 0x00
*/
func NewParser(eval string) *parser {
	tokens, err := tokenize(eval)
	if err != nil {
		return nil
	}
	return &parser{
		t:   tokens,
		idx: 0,
	}
}

func (p *parser) GetAST() *AstNode {
	return p.parseOr()
}

func (a *AstNode) Pretty(indent string, isRight bool) {
	if a == nil {
		return
	}
	// 先打印右子树
	if a.Right != nil {
		newIndent := indent
		if isRight {
			newIndent += "        "
		} else {
			newIndent += " |      "
		}
		a.Right.Pretty(newIndent, true)
	}

	// 打印当前节点
	fmt.Print(indent)
	if isRight {
		fmt.Print(" /")
	} else {
		fmt.Print(" \\")
	}
	fmt.Printf("---- (%v,%s)\n", a.Type, a.Value)

	// 打印左子树
	if a.Left != nil {
		newIndent := indent
		if isRight {
			newIndent += " |      "
		} else {
			newIndent += "        "
		}
		a.Left.Pretty(newIndent, false)
	}
}

func (p *parser) parseOr() *AstNode {
	left := p.parseAnd()

	for p.t[p.idx].toType == TokenOr {
		value := p.nowToken().t
		p.advanceToken()
		right := p.parseAnd()
		left = &AstNode{
			Type:  NodeOp,
			Value: value,
			Left:  left,
			Right: right,
		}
	}
	return left
}

func (p *parser) parseAnd() *AstNode {
	left := p.parseCompare()

	for p.t[p.idx].toType == TokenAnd {
		value := p.nowToken().t
		p.advanceToken()
		right := p.parseCompare()
		left = &AstNode{
			Type:  NodeOp,
			Value: value,
			Left:  left,
			Right: right,
		}
	}
	return left
}

func (p *parser) parseCompare() *AstNode {
	left := p.parsePrimary()

	for p.t[p.idx].toType.isCompare() {
		value := p.nowToken().t
		p.advanceToken()
		right := p.parsePrimary()
		left = &AstNode{
			Type:  NodeOp,
			Value: value,
			Left:  left,
			Right: right,
		}
	}
	return left
}

func (p *parser) parsePrimary() *AstNode {
	// 匹配括号
	if p.nowToken().toType == TokenLB {
		p.advanceToken()
		res := p.parseOr()
		p.advanceToken()
		return res
	}
	// 基本元素
	if p.nowToken().toType == TokenVar {
		value := p.nowToken().t
		p.advanceToken()
		return &AstNode{
			Type:  NodeVar,
			Value: value,
		}
	}

	if p.nowToken().toType == TokenConst {
		value := p.nowToken().t
		p.advanceToken()
		return &AstNode{
			Type:  NodeConst,
			Value: value,
		}
	}

	return nil
}

func getOpTokenType(t string) tokenType {
	switch t {
	case "&&":
		return TokenAnd
	case "||":
		return TokenOr
	case "<":
		return TokenLt
	case "<=":
		return TokenLEt
	case ">":
		return TokenMt
	case ">=":
		return TokenMEt
	case "==":
		return TokenEq
	case "!=":
		return TokenNotEq
	default:
		return -1
	}
}

func isVariable(token string) bool {
	if token[0] != '{' {
		return false
	}

	if token[len(token)-1] != '}' {
		return false
	}

	return true
}

// 十进制数和十六进制数，只判断整数
func isConstant(token string) bool {
	if _, err := strconv.ParseInt(token, 0, 64); err != nil {
		return true
	}
	return false
}

func (p *parser) nowToken() *token {
	return &p.t[p.idx]
}

func (p *parser) advanceToken() error {
	if p.idx == len(p.t)-1 {
		return errors.New("last token")
	}
	p.idx++
	return nil
}

func isHexDigit(ch byte) bool {
	return unicode.IsDigit(rune(ch)) || (ch >= 'a' && ch <= 'f') || (ch >= 'A' && ch <= 'F')
}

func isOperatorChar(ch byte) bool {
	switch ch {
	case '=', '!', '>', '<', '&', '|':
		return true
	}
	return false
}

func isValidSingleOpChar(ch byte) bool {
	switch ch {
	case '>', '<':
		return true
	}
	return false
}

func isTwoCharOp(op string) bool {
	switch op {
	case "==", "!=", ">=", "<=", "&&", "||":
		return true
	}
	return false
}
