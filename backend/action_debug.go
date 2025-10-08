package bsd_testtool

import (
	"time"
)

type PrintAction struct {
	PrintActionFeatureField
}

type DelayAction struct {
	DelayActionFeatureField
}

func (p *PrintAction) doAction(ctx *ActionContext) error {
	printFmt := p.PrintString

	printStr := FmtSprintf(printFmt, ctx)

	var err error
	if ctx.log != nil {
		_, err = ctx.log.Write([]byte(printStr))
	}

	ctx.DefaultNextAction()

	return err
}

func (d *DelayAction) doAction(ctx *ActionContext) error {
	ms := d.DelayMs
	time.Sleep(time.Duration(ms) * time.Millisecond)

	ctx.DefaultNextAction()

	return nil
}
