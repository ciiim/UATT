package bsd_testtool

type IAction interface {
	doAction(ctx *ActionContext) error
}
