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

	// For循环状态控制
	// For和If Else一样, 在条件为假的时候会直接跳转到EndBlock, 但是For循环里的语句执行完毕之后也会到EndBlock
	// 加一个状态变量，在执行For语句的时候，每次都会对这个变量赋值，代表是要跳出循环还是正常执行后回到For语句
	ForExit int
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

type AssignAction struct {
	AssignFeatureField
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

	fmt.Printf("if result:%v", res)

	resBool := res.(bool)

	if resBool {
		con.nextUID = cf.TrueUID
	} else {
		con.nextUID = cf.FalseUID
	}

	fmt.Printf("next UID:%v", con.nextUID)

	ctx.SetController(con)

	return nil
}

func (e *ElseAction) doAction(ctx *ActionContext) error {

	cf := ctx.controlFlowMap[ctx.nowAction.actionUID]

	con := &EnginControllor{nextUID: cf.EndUID}

	ctx.SetController(con)

	return nil
}

// if else for 块闭合语句
// 如果遇到for循环，要跳回for开始
func (e *EndBlockAction) doAction(ctx *ActionContext) error {
	if e.ForExit == 0 {
		cf, has := ctx.controlFlowMap[ctx.nowAction.actionUID]
		if !has {
			return fmt.Errorf("this action need a control flow")
		}
		fmt.Printf("loop start:%d\n", cf.LoopStartUID)
		con := &EnginControllor{nextUID: cf.LoopStartUID}
		ctx.SetController(con)
	} else {
		ctx.DefaultNextAction()
	}

	return nil
}

func (f *ForAction) doAction(ctx *ActionContext) error {

	// 检查跳出表达式
	forControl := ctx.controlFlowMap[ctx.nowAction.actionUID]

	con := &EnginControllor{nextUID: StopUID}

	res, err := forControl.expressionAst.Eval(ctx)
	if err != nil {
		ctx.SetController(con)
		return fmt.Errorf("eval err:%s", err.Error())
	}

	fmt.Printf("expressionAst eval res:%+v\n", res)

	resBool := false

	switch t := res.(type) {
	case int:
		if t > 0 {
			resBool = true
		} else {
			resBool = false
		}
	case bool:
		resBool = t
	}

	// 获取EndBlock
	endAction, has := ctx.actionUIDMap[forControl.EndUID]
	if !has {
		ctx.SetController(con)
		return fmt.Errorf("this for need a endblock")
	}
	ea := endAction.action.(*EndBlockAction)

	if resBool {
		con.nextUID = -1
		ea.ForExit = 0

		assignVar, has := ctx.variableMap[f.UseVar]
		if !has {
			con.nextUID = StopUID
			ctx.SetController(con)
			return fmt.Errorf("no var named:%s", f.UseVar)
		}

		forCf := ctx.controlFlowMap[ctx.nowAction.actionUID]

		res, err := forCf.assignAst.Eval(ctx)
		if err != nil {
			con.nextUID = StopUID
			ctx.SetController(con)
			return fmt.Errorf("eval failed, %s", err.Error())
		}

		assignVar.v = res

	} else {
		con.nextUID = forControl.EndUID
		ea.ForExit = 1
	}
	ctx.SetController(con)

	return nil
}

// goto标签
func (l *LabelAction) doAction(ctx *ActionContext) error {

	ctx.DefaultNextAction()

	return nil

}

// goto
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

// 更改波特率
func (c *ChangeBaudRateAction) doAction(ctx *ActionContext) error {

	ctx.DefaultNextAction()

	return nil
}

// 停止
func (s *StopAction) doAction(ctx *ActionContext) error {

	con := &EnginControllor{nextUID: StopUID}

	ctx.SetController(con)

	if s.StopCode != 0 {
		return fmt.Errorf("stop code %d", s.StopCode)
	}

	return nil
}

// 赋值
func (a *AssignAction) doAction(ctx *ActionContext) error {

	assignAst := ctx.assignMap[ctx.nowAction.actionUID]

	con := &EnginControllor{nextUID: StopUID}

	if assignAst == nil {
		ctx.SetController(con)
		return fmt.Errorf("cannot found UID:%d's ASTNode", ctx.nowAction.actionUID)
	}

	res, err := assignAst.Eval(ctx)
	if err != nil {
		ctx.SetController(con)
		return err
	}

	targetVar := ctx.variableMap[a.AssignTargetVar]
	if targetVar == nil {
		ctx.SetController(con)
		return fmt.Errorf("cannot found var:%s", a.AssignTargetVar)
	}
	targetVar.v = res

	ctx.DefaultNextAction()

	return nil
}
