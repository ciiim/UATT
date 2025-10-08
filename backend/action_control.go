package bsd_testtool

import (
	"fmt"
)

type DeclareAction struct {
	DeclareActionFeatureField
}

type IfAction struct {
	IfActionFeatureField
}

type ElseAction struct {
	ElseActionFeatureField
}

type EndBlockAction struct {
	BlockEndActionFeatureField
}

type ForAction struct {
	ForActionFeatureField
}

type LabelAction struct {
	LabelActionFeatureField
}

type GotoAction struct {
	GotoActionFeatureField
}

type ChangeBaudRateAction struct {
	ChangeBaudRateActionFeatureField
}

type StopAction struct {
	StopActionFeatureField
}

func (d *DeclareAction) doAction(ctx *ActionContext) error {

	ctx.DefaultNextAction()

	return nil
}

func (i *IfAction) doAction(ctx *ActionContext) error {

	cf := ctx.controlFlowMap[ctx.nowAction.actionUID]

	con := &EnginControllor{nextUID: StopUID}

	res, err := cf.expressionAst.Eval(ctx)
	if err != nil {
		ctx.SetController(con)
		return err
	}

	resBool := res.(bool)

	if resBool {
		con.nextUID = cf.TrueUID
	} else {
		con.nextUID = cf.FalseUID
	}

	ctx.SetController(con)

	return nil
}

func (e *ElseAction) doAction(ctx *ActionContext) error {

	ctx.DefaultNextAction()

	return nil
}

func (e *EndBlockAction) doAction(ctx *ActionContext) error {

	ctx.DefaultNextAction()

	return nil
}

func (f *ForAction) doAction(ctx *ActionContext) error {

	return nil
}

func (l *LabelAction) doAction(ctx *ActionContext) error {

	ctx.DefaultNextAction()

	return nil

}

func (g *GotoAction) doAction(ctx *ActionContext) error {

	label := g.Label

	a, has := ctx.labelMap[label]
	if !has {
		return ErrInvalidLabel
	}

	con := &EnginControllor{nextUID: a}

	ctx.SetController(con)

	return nil

}

func (c *ChangeBaudRateAction) doAction(ctx *ActionContext) error {

	ctx.DefaultNextAction()

	return nil
}

func (s *StopAction) doAction(ctx *ActionContext) error {

	con := &EnginControllor{nextUID: StopUID}

	ctx.SetController(con)

	if s.StopCode != 0 {
		return fmt.Errorf("stop code %d", s.StopCode)
	}

	return nil
}
