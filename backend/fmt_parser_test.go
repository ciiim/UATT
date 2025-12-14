package bsd_testtool_test

import (
	bsd_testtool "bsd_testtool/backend"
	"fmt"
	"testing"
)

func TestFmtParse(t *testing.T) {
	fmt.Printf("parse: %v\n", bsd_testtool.FmtSprintf("action_name {0}, recv {1:2}, res {2}, {test}",
		&bsd_testtool.TestActionContext))
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

func TestAdd(t *testing.T) {
	str := "{testNumber} + 1 + 10"

	p := bsd_testtool.NewParser(str)

	t.Logf("parser: %v\n", p)

	ast := p.GetAST()

	t.Logf("ast: %v\n", ast)

	ast.Pretty("", true)

	is, err := ast.Eval(&bsd_testtool.TestActionContext)
	if err != nil {
		t.Error(err)
	}

	fmt.Printf("is: %v\n", is)
}

func TestParser(t *testing.T) {
	str := "{1:1} == 0x01 && (1 == 1 || 1 > 0)"

	p := bsd_testtool.NewParser(str)

	t.Logf("parser: %v\n", p)

	ast := p.GetAST()

	t.Logf("ast: %v\n", ast)

	ast.Pretty("", true)

	is, err := ast.Eval(&bsd_testtool.TestActionContext)
	if err != nil {
		t.Error(err)
	}

	fmt.Printf("is: %v\n", is)

}
