package bsd_testtool

import (
	"encoding/json"
	"errors"
)

type SubModule any

type IOSubModuleBase struct {
	ModuleTypeID int `json:"ModuleTypeID"`
	ModuleUID    int `json:"ModuleUID"`
}

type IOSubModuleFill struct {
	IOSubModuleBase
	FillLength int    `json:"FillLength"`
	UseVar     string `json:"UseVar"`
}

type IOSubModuleFixed struct {
	IOSubModuleBase
	FixedContent []int `json:"FixedContent"`
}

type IOSubModuleCalc struct {
	IOSubModuleBase
	Mode                string `json:"Mode"`
	CalcFunc            string `json:"CalcFunc"`
	CalcOrder           int    `json:"CalcOrder"`
	CalcTiming          string `json:"CalcTiming"`
	PlaceholderBytes    []int  `json:"PlaceholderBytes"`
	CalcInputModulesUID []int  `json:"CalcInputModulesUID"`
}

type IOSubModuleCustom struct {
	IOSubModuleBase
	CustomLength int `json:"CustomLength"`
}

type IOModuleFeatureField struct {
	TimeoutMs  int
	SubModules []SubModule
}

func (i *IOModuleFeatureField) UnmarshalJSON(b []byte) error {
	type Temp struct {
		TimeoutMs  int               `json:"TimeoutMs"`
		SubModules []json.RawMessage `json:"SubModules"`
	}

	var t Temp

	if err := json.Unmarshal(b, &t); err != nil {
		return err
	}

	i.TimeoutMs = t.TimeoutMs

	for _, raw := range t.SubModules {
		var base IOSubModuleBase
		if err := json.Unmarshal(raw, &base); err != nil {
			return err
		}
		switch SubModuleTypeID(base.ModuleTypeID) {
		case FillSubMT:
			var fill IOSubModuleFill
			if err := json.Unmarshal(raw, &fill); err != nil {
				return err
			}
			i.SubModules = append(i.SubModules, &fill)
		case FixedSubMT:
			var fixed IOSubModuleFixed
			if err := json.Unmarshal(raw, &fixed); err != nil {
				return err
			}
			i.SubModules = append(i.SubModules, &fixed)
		case CalcSubMT:
			var calc IOSubModuleCalc
			if err := json.Unmarshal(raw, &calc); err != nil {
				return err
			}
			i.SubModules = append(i.SubModules, &calc)
		case CustomSubMT:
			var custom IOSubModuleCustom
			if err := json.Unmarshal(raw, &custom); err != nil {
				return err
			}
			i.SubModules = append(i.SubModules, &custom)
		default:
			return errors.New("unsupport sub module")
		}
	}

	return nil
}
