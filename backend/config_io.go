package bsd_testtool

import (
	"bsd_testtool/backend/types"
	"encoding/json"
	"errors"
)

type IConfigModule interface {
	ToModule() IOModule
}

type IOModuleConfigBase struct {
	Index        int
	ModuleTypeID int `json:"ModuleTypeID"`
	ModuleUID    int `json:"ModuleUID"`
}

type IOModuleConfigFill struct {
	IOModuleConfigBase
	FillLength int    `json:"FillLength"`
	UseVar     string `json:"UseVar"`
}

func (f *IOModuleConfigFill) ToModule() IOModule {
	return &IOFillModule{*f}
}

type IOModuleConfigFixed struct {
	IOModuleConfigBase
	FixedContent []int `json:"FixedContent"`
}

func (f *IOModuleConfigFixed) ToModule() IOModule {
	return &IOFixedModule{*f}
}

type IOModuleConfigCalc struct {
	IOModuleConfigBase
	Mode                string `json:"Mode"`
	CalcFunc            string `json:"CalcFunc"`
	CalcTiming          string `json:"CalcTiming"`
	PlaceholderBytes    []int  `json:"PlaceholderBytes"`
	CalcInputModulesUID []int  `json:"CalcInputModulesUID"`
}

func (c *IOModuleConfigCalc) ToModule() IOModule {
	return &IOCalcModule{*c}
}

type IOModuleConfigCustom struct {
	IOModuleConfigBase
	// 只有接收会用到
	// 该模块的长度由指定UID的模块接收到的内容决定
	// 此时检查操作不会起作用
	ReceiveVarLengthModuleUID types.ModuleUID `json:"ReceiveVarLengthModuleUID"`
	CustomContent             []int           `json:"CustomContent"`
}

func (c *IOModuleConfigCustom) ToModule() IOModule {
	return &IOCustomModule{*c}
}

type IOActionFeatureField struct {
	TimeoutMs int
	Modules   []IConfigModule
}

type IOSendActionFeatureField struct {
	IOActionFeatureField
}

func (s *IOSendActionFeatureField) ToAction() IAction {
	modules := make([]IOModule, len(s.Modules))
	for i, cm := range s.Modules {
		modules[i] = cm.ToModule()
	}
	return &SendAction{
		IOAction: IOAction{
			TimeoutMs: s.TimeoutMs,
			Modules:   modules,
		},
	}
}

type IOReceiveActionFeatureField struct {
	IOActionFeatureField
}

func (r *IOReceiveActionFeatureField) ToAction() IAction {
	modules := make([]IOModule, len(r.Modules))
	for i, cm := range r.Modules {
		modules[i] = cm.ToModule()
	}
	return &ReceiveAction{
		IOAction: IOAction{
			TimeoutMs: r.TimeoutMs,
			Modules:   modules,
		},
	}
}

func unmarshalIOAction(actionTypeID types.ActionTypeID, b []byte) (any, error) {
	switch actionTypeID {
	case types.SendAT:
		var s IOSendActionFeatureField
		if err := json.Unmarshal(b, &s); err != nil {
			return nil, err
		}
		return &s, nil
	case types.ReceiveAT:
		var r IOReceiveActionFeatureField
		if err := json.Unmarshal(b, &r); err != nil {
			return nil, err
		}
		return &r, nil
	default:
		return nil, errors.New("unsupport action")
	}
}

func (i *IOActionFeatureField) UnmarshalJSON(b []byte) error {
	type Temp struct {
		TimeoutMs int               `json:"TimeoutMs"`
		Modules   []json.RawMessage `json:"Modules"`
	}

	var t Temp

	if err := json.Unmarshal(b, &t); err != nil {
		return err
	}

	i.TimeoutMs = t.TimeoutMs

	for _, raw := range t.Modules {
		var base IOModuleConfigBase
		if err := json.Unmarshal(raw, &base); err != nil {
			return err
		}
		switch types.ModuleTypeID(base.ModuleTypeID) {
		case types.FillMT:
			var fill IOModuleConfigFill
			if err := json.Unmarshal(raw, &fill); err != nil {
				return err
			}
			i.Modules = append(i.Modules, &fill)
		case types.FixedMT:
			var fixed IOModuleConfigFixed
			if err := json.Unmarshal(raw, &fixed); err != nil {
				return err
			}
			i.Modules = append(i.Modules, &fixed)
		case types.CalcMT:
			var calc IOModuleConfigCalc
			if err := json.Unmarshal(raw, &calc); err != nil {
				return err
			}
			i.Modules = append(i.Modules, &calc)
		case types.CustomMT:
			var custom IOModuleConfigCustom
			if err := json.Unmarshal(raw, &custom); err != nil {
				return err
			}
			i.Modules = append(i.Modules, &custom)
		default:
			return errors.New("unsupport sub module")
		}
	}

	return nil
}
