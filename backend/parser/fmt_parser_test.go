package parser_test

import (
	"bsd_testtool/backend/context"
	"bsd_testtool/backend/parser"
	"errors"
	"fmt"
	"testing"
)

func TestFmtParse(t *testing.T) {
	var parser parser.StringFmt
	fmt.Printf("parse: %v\n", parser.Sprintf("action_name {0}, recv {1}, res {2}", &context.Context{
		LastModuleName:   "Send 11",
		LastSerialBuffer: []byte{0xAA, 0xFF, 0x01, 0x05},
		LastExecResult:   errors.New("Test Error"),
	}))
}
