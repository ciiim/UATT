package bsd_testtool

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode"
)

var testVar variable = variable{
	varName: "test",
	varType: VarNumberArray,
	v:       []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 0},
}

var testNumberVar variable = variable{
	varName: "testNumber",
	varType: VarNumber,
	v:       int(100),
}

var TestActionContext ActionContext = ActionContext{
	LastActionName:   "Send 11",
	LastSerialBuffer: []byte{0xAA, 0xFF, 0x01, 0x05, 0x55, 0x55, 0x55},
	LastExecResult:   errors.New("test error"),
	variableMap: map[string]*variable{
		"test":       &testVar,
		"testNumber": &testNumberVar,
	},
}

func FmtSprintf(format string, ctx *ActionContext) string {
	re := regexp.MustCompile(`\{(\S+)\}`)
	return re.ReplaceAllStringFunc(format, func(s string) string {
		if ctx == nil {
			return "[ctx nil]"
		}

		v := re.FindStringSubmatch(s)[1]

		res := FmtGetVar(v, ctx)
		if res == nil {
			return "[nil]"
		}
		switch resVar := res.(type) {
		case string:
			return resVar
		default:
			return fmt.Sprintf("%v", resVar)
		}
	})
}

// 取变量值，如{test}, {res:1,5}, 输入varName字符串不带{}
// 返回值为整数或整数数组
func FmtGetVar(varFmt string, ctx *ActionContext) any {
	arrFmtStrings := strings.Split(varFmt, ":")
	arrIndexBegin, arrIndexEnd := -1, -1
	varName := arrFmtStrings[0]

	if len(arrFmtStrings) > 1 {
		arrIndexStrings := strings.Split(arrFmtStrings[1], ",")

		if len(arrIndexStrings) == 2 { // 范围
			end, err := strconv.Atoi(arrIndexStrings[1])
			if err == nil {
				arrIndexEnd = end
			}
		}
		begin, err := strconv.Atoi(arrIndexStrings[0])
		if err == nil {
			arrIndexBegin = begin
		}

		// 解析错误
		if arrIndexBegin == -1 {
			return nil
		}

		// 上界小于下界
		if arrIndexEnd != -1 && arrIndexEnd < arrIndexBegin {
			return nil
		}
	}
	switch varName[0] {
	case '0':
		return ctx.LastActionName
	case '1':
		if arrIndexBegin > len(ctx.LastSerialBuffer)-1 || arrIndexEnd > len(ctx.LastSerialBuffer)-1 {
			return nil
		}
		if arrIndexBegin != -1 && arrIndexEnd == -1 {
			return ctx.LastSerialBuffer[arrIndexBegin]
		} else if arrIndexBegin != -1 {
			return ctx.LastSerialBuffer[arrIndexBegin:arrIndexEnd]
		} else {
			return ctx.LastSerialBuffer
		}
	case '2':
		return func() any {
			if ctx.LastExecResult == nil {
				return nil
			} else {
				return ctx.LastExecResult.Error()
			}
		}()
	case '3':
		return time.Now().Format("01-02 15:04:05")
	default:
		v, has := ctx.variableMap[varName]
		if !has {
			return nil
		}

		var vArr []int = nil

		switch v.varType {
		case VarNumber:
			vArr = make([]int, 1)
			vArr[0] = v.v.(int)
		case VarNumberArray:
			vArr = v.v.([]int)
		case VarString:
			runes := []rune(v.v.(string))
			vArr = make([]int, len(runes))
			for i, r := range runes {
				vArr[i] = int(r)
			}
		case VarJSON:
			runes := []rune(v.v.(string))
			vArr = make([]int, len(runes))
			for i, r := range runes {
				vArr[i] = int(r)
			}
		default:
			return nil
		}

		if arrIndexEnd > len(vArr)-1 {
			return nil
		}

		if arrIndexBegin > len(vArr)-1 {
			return nil
		}

		if arrIndexEnd != -1 {
			return vArr[arrIndexBegin:arrIndexEnd]
		} else if arrIndexBegin != -1 {
			return vArr[arrIndexBegin]
		} else {
			if v.varType == VarNumber {
				return vArr[0]
			} else {
				return vArr
			}
		}
	}

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
		// {0} 上下文里的 LastActionName
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

	TokenAdd   // +
	TokenMinus // -
)

func (t tokenType) isCompare() bool {
	return (t >= TokenLt && t <= TokenNotEq)
}

func (t tokenType) isAddMinus() bool {
	return (t >= TokenAdd && t <= TokenMinus)
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

func compareEq(lv, rv int) bool {
	return lv == rv
}

func compareNum(lv, rv int) int {
	return lv - rv
}

func toBool(v int) bool {
	return v > 0
}

func addOp(lv, rv int) int {
	return lv + rv
}

func minusOp(lv, rv int) int {
	return lv - rv
}

func (a *AstNode) Eval(ctx *ActionContext) (any, error) {
	switch a.Type {
	case NodeOp:
		lv, err := a.Left.Eval(ctx)
		if err != nil {
			return false, err
		}

		rv, err := a.Right.Eval(ctx)
		if err != nil {
			return false, err
		}

		lvo := 0
		rvo := 0

		// 把所有类型的值都转换成int
		switch t := lv.(type) {
		case int:
			lvo = t
		case bool:
			lvo = func() int {
				if t {
					return 1
				}
				return 0
			}()
		case int64:
			lvo = int(t)
		case uint8:
			lvo = int(t)
		default:
			return nil, fmt.Errorf("lv unsupport type %V", t)
		}

		switch t := rv.(type) {
		case int:
			rvo = t
		case bool:
			rvo = func() int {
				if t {
					return 1
				}
				return 0
			}()
		case int64:
			rvo = int(t)
		case uint8:
			rvo = int(t)
		default:
			return nil, fmt.Errorf("rv unsupport type %V", t)
		}

		// 根据运算符做实际运算
		switch a.Value {
		case "==":
			return compareEq(lvo, rvo), nil
		case "!=":
			return !compareEq(lvo, rvo), nil
		case "||":
			return toBool(lvo) || toBool(rvo), nil
		case "&&":
			return toBool(lvo) && toBool(rvo), nil
		case "<":
			return compareNum(lvo, rvo) < 0, nil
		case ">":
			return compareNum(lvo, rvo) > 0, nil
		case "<=":
			return compareNum(lvo, rvo) <= 0, nil
		case ">=":
			return compareNum(lvo, rvo) >= 0, nil
		case "+":
			return addOp(lvo, rvo), nil
		case "-":
			return minusOp(lvo, rvo), nil
		}
	case NodeConst:
		n, err := strconv.ParseInt(a.Value, 0, 64)
		if err != nil {
			return false, err
		}
		return n, nil
	case NodeVar:
		// FmtGetVar要求输入varFmt不带{}
		v := FmtGetVar(a.Value[1:len(a.Value)-1], ctx)
		switch t := v.(type) {
		case string:
			return false, errors.New("node var cannot be string")
		case nil:
			return nil, fmt.Errorf("wrong var[%v]", a.Value[1:len(a.Value)-1])
		default:
			return t, nil
		}
	}
	return false, errors.New("empty ast")
}

/*
优先级
低到高
0. ||
1. &&
2. ==, !=, >, <, >=, <=
3. +, -
5. 括号

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
	if len(p.t) == 0 {
		return nil
	}
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
		_ = p.advanceToken()
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
		_ = p.advanceToken()
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
	left := p.parseAddMinus()

	for p.t[p.idx].toType.isCompare() {
		value := p.nowToken().t
		_ = p.advanceToken()
		right := p.parseAddMinus()
		left = &AstNode{
			Type:  NodeOp,
			Value: value,
			Left:  left,
			Right: right,
		}
	}
	return left
}

func (p *parser) parseAddMinus() *AstNode {
	left := p.parsePrimary()

	for p.t[p.idx].toType.isAddMinus() {
		value := p.nowToken().t
		_ = p.advanceToken()
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
		_ = p.advanceToken()
		res := p.parseOr()
		_ = p.advanceToken()
		return res
	}
	// 基本元素
	if p.nowToken().toType == TokenVar {
		value := p.nowToken().t
		_ = p.advanceToken()
		return &AstNode{
			Type:  NodeVar,
			Value: value,
		}
	}

	if p.nowToken().toType == TokenConst {
		value := p.nowToken().t
		_ = p.advanceToken()
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
	case "+":
		return TokenAdd
	case "-":
		return TokenMinus
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
	case '=', '!', '>', '<', '&', '|', '+', '-':
		return true
	}
	return false
}

func isValidSingleOpChar(ch byte) bool {
	switch ch {
	case '>', '<', '+', '-':
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
