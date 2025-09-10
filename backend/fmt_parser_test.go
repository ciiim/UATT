package bsd_testtool_test

import (
	bsd_testtool "bsd_testtool/backend"
	"errors"
	"fmt"
	"testing"
)

func TestFmtParse(t *testing.T) {
	var parser bsd_testtool.StringFmt
	fmt.Printf("parse: %v\n", parser.Sprintf("action_name {0}, recv {1}, res {2}", &bsd_testtool.ActionContext{
		LastModuleName:   "Send 11",
		LastSerialBuffer: []byte{0xAA, 0xFF, 0x01, 0x05},
		LastExecResult:   errors.New("Test Error"),
	}))
}
