package bsd_testtool

import (
	"bsd_testtool/backend/types"
	"encoding/json"
	"errors"
)

type PrintActionFeatureField struct {
	PrintString string `json:"PrintFmt"`
}

func (p *PrintActionFeatureField) ToAction() IAction {
	return &PrintAction{*p}
}

type DelayActionFeatureField struct {
	DelayMs int `json:"DelayMs"`
}

func (d *DelayActionFeatureField) ToAction() IAction {
	return &DelayAction{*d}
}

type ShowActionFeatureField struct {
	OutputID int    `json:"OutputID"`
	FmtStr   string `json:"FmtStr"`
}

func (s *ShowActionFeatureField) ToAction() IAction {
	return &ShowAction{*s}
}

func unmarshalDebugAction(actionTypeID types.ActionTypeID, b []byte) (any, error) {
	switch actionTypeID {
	case types.PrintAT:
		var p PrintActionFeatureField
		if err := json.Unmarshal(b, &p); err != nil {
			return nil, err
		}
		return &p, nil
	case types.DelayAT:
		var d DelayActionFeatureField
		if err := json.Unmarshal(b, &d); err != nil {
			return nil, err
		}
		return &d, nil
	case types.ShowAT:
		var s ShowActionFeatureField
		if err := json.Unmarshal(b, &s); err != nil {
			return nil, err
		}
		return &s, nil
	default:
		return nil, errors.New("unsupport action")
	}
}
