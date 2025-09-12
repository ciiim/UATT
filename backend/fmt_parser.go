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

type tokenType int

var ErrVarNotClosed error = errors.New("var not closed")
var ErrInvalidBinaryOp error = errors.New("invalid binary operator")
var ErrInvalidToken error = errors.New("invalid token")

func TestTokenize(input string) ([]string, error) {
	return tokenize(input)
}

func tokenize(input string) ([]string, error) {
	var tokens []string
	i := 0
	for i < len(input) {
		ch := input[i]

		// 跳过空格
		if unicode.IsSpace(rune(ch)) {
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
				tokens = append(tokens, input[i:j+1])
				fmt.Printf("var\n")
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

			tokens = append(tokens, input[i:j])
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
					tokens = append(tokens, op)
					i += 2
					continue
				}
			}
			if !isValidSingleOpChar(ch) {
				return nil, ErrInvalidToken
			}
			// 单字符运算符
			tokens = append(tokens, string(ch))
			i++
			continue
		}

		return nil, ErrInvalidToken

	}

	return tokens, nil
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
