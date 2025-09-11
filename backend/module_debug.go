package bsd_testtool

import (
	"encoding/json"
	"errors"
	"time"
)

type PrintModuleFeatureField struct {
	PrintString string `json:"PrintFmt"`
}

type DelayModuleFeatureField struct {
	DelayMs int `json:"DelayMs"`
}

func unmarshalDebugModule(moduleTypeID ModuleTypeID, b []byte) (any, error) {
	switch moduleTypeID {
	case PrintMT:
		var p PrintModuleFeatureField
		if err := json.Unmarshal(b, &p); err != nil {
			return nil, err
		}
		return &p, nil
	case DelayMT:
		var d DelayModuleFeatureField
		if err := json.Unmarshal(b, &d); err != nil {
			return nil, err
		}
		return &d, nil
	default:
		return nil, errors.New("unsupport module")
	}
}

func doPrint(ctx *ActionContext, m *ModuleBase) error {
	printFmt := m.TypeFeatureField.(PrintModuleFeatureField).PrintString

	printStr := FmtSprintf(printFmt, ctx)

	_, err := ctx.log.Write([]byte(printStr))

	return err
}

func doDelay(ctx *ActionContext, m *ModuleBase) error {
	ms := m.TypeFeatureField.(DelayModuleFeatureField).DelayMs
	time.Sleep(time.Duration(ms) * time.Millisecond)
	return nil
}
