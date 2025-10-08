package bsd_testtool_test

import (
	bsd_testtool "bsd_testtool/backend"
	"errors"
	"fmt"
	"testing"
)

func TestFmtParse(t *testing.T) {
	fmt.Printf("parse: %v\n", bsd_testtool.FmtSprintf("action_name {0}, recv {1}, res {2}", &bsd_testtool.ActionContext{
		LastActionName:   "Send 11",
		LastSerialBuffer: []byte{0xAA, 0xFF, 0x01, 0x05},
		LastExecResult:   errors.New("Test Error"),
	}))
}

func TestTokenize(t *testing.T) {
	res, err := bsd_testtool.TestTokenize("({varA}== 0x01 || {varB} > 1) && 1 != 1")
	if err != nil {
		t.Error(err)
		return
	}

	for _, token := range res {
		t.Logf("token:[%v]\n", token)
	}
}

func TestParser(t *testing.T) {
	str := "{1:6} == 1"

	p := bsd_testtool.NewParser(str)

	t.Logf("parser: %v\n", p)

	ast := p.GetAST()

	t.Logf("ast: %v\n", ast)

	ast.Pretty("", true)

}
