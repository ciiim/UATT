package bsd_testtool

import (
	"encoding/json"
	"errors"
)

type DeclareModuleFeatureField struct {
	VarName           string `json:"VarName"`
	VarType           string `json:"VarType"`
	VarNumberValue    int    `json:"VarNumberValue,omitempty"`
	VarStringValue    string `json:"VarStringValue,omitempty"`
	VarByteArrayValue []byte `json:"VarByteArrayValue,omitempty"`
}

type IfModuleFeatureField struct {
	Condition string `json:"Condition"`
}

type ElseModuleFeatureField struct {
}

type ForModuleFeatureField struct {
	UseVar         string `json:"UseVar"`
	EnterCondition string `json:"EnterCondition"`
	VarOp          string `json:"VarOp"`
}

type BlockEndModuleFeatureField struct {
}

type LabelModuleFeatureField struct {
	LabelName string `json:"LabelName"`
}

type GotoModuleFeatureField struct {
	Label string `json:"Label"`
}

type ChangeBaudRateModuleFeatureField struct {
	TargetBaudRate int `json:"TargetBaudRate"`
}

type StopModuleFeatureField struct {
	StopCode int `json:"StopCode"`
}

func unmarshalControlModule(moduleTypeID ModuleTypeID, b []byte) (any, error) {
	switch moduleTypeID {
	case IfMT:
		var f IfModuleFeatureField
		if err := json.Unmarshal(b, &f); err != nil {
			return nil, err
		}
		return &f, nil
	case ElseMT:
		var f ElseModuleFeatureField
		if err := json.Unmarshal(b, &f); err != nil {
			return nil, err
		}
		return &f, nil
	case ForLabelMT:
		var f ForModuleFeatureField
		if err := json.Unmarshal(b, &f); err != nil {
			return nil, err
		}
		return &f, nil
	case EndBlockMT:
		var f BlockEndModuleFeatureField
		if err := json.Unmarshal(b, &f); err != nil {
			return nil, err
		}
		return &f, nil
	case GotoLabelMT:
		var f LabelModuleFeatureField
		if err := json.Unmarshal(b, &f); err != nil {
			return nil, err
		}
		return &f, nil
	case GotoMT:
		var f GotoModuleFeatureField
		if err := json.Unmarshal(b, &f); err != nil {
			return nil, err
		}
		return &f, nil
	case ChangeBaudRateMT:
		var f ChangeBaudRateModuleFeatureField
		if err := json.Unmarshal(b, &f); err != nil {
			return nil, err
		}
		return &f, nil
	case StopMT:
		var f StopModuleFeatureField
		if err := json.Unmarshal(b, &f); err != nil {
			return nil, err
		}
		return &f, nil

	default:
		return nil, errors.New("unsupport module")
	}
}

func doIf(ctx *ActionContext, m *ModuleBase) error {
	return nil
}

func doElse(ctx *ActionContext, m *ModuleBase) error {
	return nil

}

func doEndBlock(ctx *ActionContext, m *ModuleBase) error {
	return nil

}

func doFor(ctx *ActionContext, m *ModuleBase) error {
	return nil

}

func doLabel(ctx *ActionContext, m *ModuleBase) error {
	return nil

}

func doGoto(ctx *ActionContext, m *ModuleBase) error {
	return nil

}

func doChangeBaudRate(ctx *ActionContext, m *ModuleBase) error {
	return nil
}

func doStop(ctx *ActionContext, m *ModuleBase) error {
	return nil
}
