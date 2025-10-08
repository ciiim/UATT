package bsd_testtool

import (
	"bsd_testtool/backend/types"
	"encoding/json"
	"errors"
)

type DeclareActionFeatureField struct {
	VarName           string `json:"VarName"`
	VarType           string `json:"VarType"`
	VarNumberValue    int    `json:"VarNumberValue,omitempty"`
	VarStringValue    string `json:"VarStringValue,omitempty"`
	VarByteArrayValue []byte `json:"VarByteArrayValue,omitempty"`
}

func (d *DeclareActionFeatureField) ToAction() IAction {
	return &DeclareAction{*d}
}

type IfActionFeatureField struct {
	Condition string `json:"Condition"`
}

func (i *IfActionFeatureField) ToAction() IAction {
	return &IfAction{*i}
}

type ElseActionFeatureField struct {
}

func (e *ElseActionFeatureField) ToAction() IAction {
	return &ElseAction{*e}
}

type ForActionFeatureField struct {
	UseVar         string `json:"UseVar"`
	EnterCondition string `json:"EnterCondition"`
	VarOp          string `json:"VarOp"`
}

func (f *ForActionFeatureField) ToAction() IAction {
	return &ForAction{*f}
}

type BlockEndActionFeatureField struct {
}

func (b *BlockEndActionFeatureField) ToAction() IAction {
	return &EndBlockAction{*b}
}

type LabelActionFeatureField struct {
	LabelName string `json:"LabelName"`
}

func (l *LabelActionFeatureField) ToAction() IAction {
	return &LabelAction{*l}
}

type GotoActionFeatureField struct {
	Label string `json:"Label"`
}

func (g *GotoActionFeatureField) ToAction() IAction {
	return &GotoAction{*g}
}

type ChangeBaudRateActionFeatureField struct {
	TargetBaudRate int `json:"TargetBaudRate"`
}

func (c *ChangeBaudRateActionFeatureField) ToAction() IAction {
	return &ChangeBaudRateAction{*c}
}

type StopActionFeatureField struct {
	StopCode int `json:"StopCode"`
}

func (s *StopActionFeatureField) ToAction() IAction {
	return &StopAction{*s}
}

func unmarshalControlAction(actionTypeID types.ActionTypeID, b []byte) (any, error) {
	switch actionTypeID {
	case types.IfAT:
		var f IfActionFeatureField
		if err := json.Unmarshal(b, &f); err != nil {
			return nil, err
		}
		return &f, nil
	case types.ElseAT:
		var f ElseActionFeatureField
		if err := json.Unmarshal(b, &f); err != nil {
			return nil, err
		}
		return &f, nil
	case types.ForLabelAT:
		var f ForActionFeatureField
		if err := json.Unmarshal(b, &f); err != nil {
			return nil, err
		}
		return &f, nil
	case types.EndBlockAT:
		var f BlockEndActionFeatureField
		if err := json.Unmarshal(b, &f); err != nil {
			return nil, err
		}
		return &f, nil
	case types.GotoLabelAT:
		var f LabelActionFeatureField
		if err := json.Unmarshal(b, &f); err != nil {
			return nil, err
		}
		return &f, nil
	case types.GotoAT:
		var f GotoActionFeatureField
		if err := json.Unmarshal(b, &f); err != nil {
			return nil, err
		}
		return &f, nil
	case types.ChangeBaudRateAT:
		var f ChangeBaudRateActionFeatureField
		if err := json.Unmarshal(b, &f); err != nil {
			return nil, err
		}
		return &f, nil
	case types.StopAT:
		var f StopActionFeatureField
		if err := json.Unmarshal(b, &f); err != nil {
			return nil, err
		}
		return &f, nil

	default:
		return nil, errors.New("unsupport action")
	}
}
