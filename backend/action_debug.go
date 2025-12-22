package bsd_testtool

import (
	"time"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type PrintAction struct {
	PrintActionFeatureField
}

type DelayAction struct {
	DelayActionFeatureField
}

type ShowAction struct {
	ShowActionFeatureField
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

func (s *ShowAction) doAction(ctx *ActionContext) error {

	showStr := FmtSprintf(s.FmtStr, ctx)

	runtime.EventsEmit(ctx.engine.wailsCtx, "output_text", s.OutputIdx, showStr)

	ctx.DefaultNextAction()

	return nil
}
