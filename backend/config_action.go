package bsd_testtool

import (
	"bsd_testtool/backend/types"
	"encoding/json"
	"errors"
)

type ActionTypeFeatureFieldBase any

type ConfigActionBase struct {

	// 模块实例UID
	ActionUID int
	// 模块类型
	ActionType string
	// 模块类型ID
	ActionTypeID int
	// 模块名 自定义
	Name string
	// 是否是断点
	BreakPoint bool
	// 模块独有属性
	TypeFeatureField ActionTypeFeatureFieldBase
}

type ConfigActionBaseJson struct {
	// 模块实例UID
	ActionUID int `json:"ActionUID"`
	// 模块类型
	ActionType string `json:"ActionType"`
	// 模块类型ID
	ActionTypeID int `json:"ActionTypeID"`
	// 模块名 自定义
	Name string `json:"Name"`
	// 是否是断点
	BreakPoint bool `json:"BreakPoint"`

	TypeFeatureField json.RawMessage `json:"TypeFeatureField"`
}

func (c ConfigActionBaseJson) ToBase() ConfigActionBase {
	base := ConfigActionBase{
		ActionUID:    c.ActionUID,
		ActionType:   c.ActionType,
		ActionTypeID: c.ActionTypeID,
		Name:         c.Name,
		BreakPoint:   c.BreakPoint,
	}

	switch types.ActionType(base.ActionType) {
	case types.IOAction:
		res, err := unmarshalIOAction(types.ActionTypeID(c.ActionTypeID), c.TypeFeatureField)
		if err != nil {
			return base
		}
		base.TypeFeatureField = res
	case types.ControlAction:
		res, err := unmarshalControlAction(types.ActionTypeID(c.ActionTypeID), c.TypeFeatureField)
		if err != nil {
			return base
		}
		base.TypeFeatureField = res
	case types.DebugAction:
		res, err := unmarshalDebugAction(types.ActionTypeID(c.ActionTypeID), c.TypeFeatureField)
		if err != nil {
			return base
		}
		base.TypeFeatureField = res
	default:
		return base
	}
	return base
}

func (ca *ConfigActionBase) UnmarshalJSON(b []byte) error {

	var t ConfigActionBaseJson

	if err := json.Unmarshal(b, &t); err != nil {
		return err
	}

	ca.ActionUID = t.ActionUID
	ca.ActionType = t.ActionType
	ca.ActionTypeID = t.ActionTypeID
	ca.Name = t.Name
	ca.BreakPoint = t.BreakPoint

	switch types.ActionType(ca.ActionType) {
	case types.IOAction:
		res, err := unmarshalIOAction(types.ActionTypeID(t.ActionTypeID), t.TypeFeatureField)
		if err != nil {
			return err
		}
		ca.TypeFeatureField = res
	case types.ControlAction:
		res, err := unmarshalControlAction(types.ActionTypeID(t.ActionTypeID), t.TypeFeatureField)
		if err != nil {
			return err
		}
		ca.TypeFeatureField = res
	case types.DebugAction:
		res, err := unmarshalDebugAction(types.ActionTypeID(t.ActionTypeID), t.TypeFeatureField)
		if err != nil {
			return err
		}
		ca.TypeFeatureField = res
	default:
		return errors.New("no support action type")
	}

	return nil
}
