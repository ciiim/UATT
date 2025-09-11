package bsd_testtool

import (
	"encoding/json"
	"errors"
)

type ModuleFunc func(ctx *ActionContext, m *ModuleBase) error

type ModuleTypeID int

type ModuleType string

type SubModuleTypeID int

const (
	IOModule      ModuleType = "IO"
	ControlModule ModuleType = "Control"
	DebugModule   ModuleType = "Debug"
)

const (
	PrintMT ModuleTypeID = 90
	DelayMT ModuleTypeID = 91

	SendMT    ModuleTypeID = 1
	ReceiveMT ModuleTypeID = 2

	DeclareMT ModuleTypeID = 23

	IfMT   ModuleTypeID = 24
	ElseMT ModuleTypeID = 25

	ForLabelMT ModuleTypeID = 26

	EndBlockMT ModuleTypeID = 27

	GotoLabelMT ModuleTypeID = 28
	GotoMT      ModuleTypeID = 29

	ChangeBaudRateMT ModuleTypeID = 30

	StopMT ModuleTypeID = 31
)

const (
	FillSubMT   SubModuleTypeID = 10
	FixedSubMT  SubModuleTypeID = 11
	CalcSubMT   SubModuleTypeID = 12
	CustomSubMT SubModuleTypeID = 13
)

type ModuleTypeFeatureFieldBase any

type ModuleBase struct {
	// 在数组内的下标
	Index int

	// 模块实例UID
	ModuleUID int
	// 模块类型
	ModuleType string
	// 模块类型ID
	ModuleTypeID int
	// 模块名 自定义
	Name string
	// 是否是断点
	BreakPoint bool
	// 模块独有属性
	TypeFeatureField ModuleTypeFeatureFieldBase
}

var ModuleFuncMap map[ModuleTypeID]ModuleFunc = map[ModuleTypeID]ModuleFunc{
	PrintMT: doPrint,
	DelayMT: doDelay,

	SendMT:    doSend,
	ReceiveMT: doReceive,

	IfMT:             doIf,
	ElseMT:           doElse,
	ForLabelMT:       doFor,
	EndBlockMT:       doEndBlock,
	GotoLabelMT:      doLabel,
	GotoMT:           doGoto,
	ChangeBaudRateMT: doChangeBaudRate,
	StopMT:           doStop,
}

func (mb *ModuleBase) UnmarshalJSON(b []byte) error {
	type Temp struct {
		// 模块实例UID
		ModuleUID int `json:"ModuleUID"`
		// 模块类型
		ModuleType string `json:"ModuleType"`
		// 模块类型ID
		ModuleTypeID int `json:"ModuleTypeID"`
		// 模块名 自定义
		Name string `json:"Name"`
		// 是否是断点
		BreakPoint bool `json:"BreakPoint"`

		TypeFeatureField json.RawMessage `json:"TypeFeatureField"`
	}

	var t Temp

	if err := json.Unmarshal(b, &t); err != nil {
		return err
	}

	mb.ModuleUID = t.ModuleUID
	mb.ModuleType = t.ModuleType
	mb.ModuleTypeID = t.ModuleTypeID
	mb.Name = t.Name
	mb.BreakPoint = t.BreakPoint

	switch ModuleType(mb.ModuleType) {
	case IOModule:
		var ioFeat IOModuleFeatureField
		if err := json.Unmarshal(t.TypeFeatureField, &ioFeat); err != nil {
			return err
		}
		mb.TypeFeatureField = &ioFeat
	case ControlModule:
		res, err := unmarshalControlModule(ModuleTypeID(t.ModuleTypeID), t.TypeFeatureField)
		if err != nil {
			return err
		}
		mb.TypeFeatureField = res
	case DebugModule:
		res, err := unmarshalDebugModule(ModuleTypeID(t.ModuleTypeID), t.TypeFeatureField)
		if err != nil {
			return err
		}
		mb.TypeFeatureField = res
	default:
		return errors.New("no support module type")
	}

	return nil
}
