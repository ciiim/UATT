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
		for _, l := range ctx.log {
			_, err = l.Write([]byte(printStr))
		}
	}

	ctx.DefaultNextAction()

	return err
}

func (d *DelayAction) doAction(ctx *ActionContext) error {

	println("delay", d.DelayMs)

	ms := d.DelayMs
	time.Sleep(time.Duration(ms) * time.Millisecond)

	ctx.DefaultNextAction()

	return nil
}
