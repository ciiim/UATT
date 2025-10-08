package bsd_testtool

import (
	"encoding/json"
	"fmt"
)

type IConfig interface {
	ToAction() IAction
}

type AppConfigSettings struct {
	AppName      string `json:"AppName"`
	SerialConfig struct {
		BaudRate int    `json:"BaudRate"`
		DataBits int    `json:"DataBits"`
		Parity   string `json:"Parity"`
		StopBits int    `json:"StopBits"`
	} `json:"SerialConfig"`
	LogEnable         bool   `json:"LogEnable"`
	LogExportEnable   bool   `json:"LogExportEnable"`
	LogExportLoaction string `json:"LogExportLoaction"`
}

type AppConfig struct {
	AppName      string `json:"AppName"`
	SerialConfig struct {
		BaudRate int    `json:"BaudRate"`
		DataBits int    `json:"DataBits"`
		Parity   string `json:"Parity"`
		StopBits int    `json:"StopBits"`
	} `json:"SerialConfig"`
	LogEnable         bool               `json:"LogEnable"`
	LogExportEnable   bool               `json:"LogExportEnable"`
	LogExportLoaction string             `json:"LogExportLoaction"`
	Actions           []ConfigActionBase `json:"Actions"`
}

func (a *AppConfig) PrintConfig() {
	fmt.Printf("app: %v\n", a)
	for _, mod := range a.Actions {
		switch m := mod.TypeFeatureField.(type) {
		case *IOSendActionFeatureField:
			fmt.Printf("IOModuleFeatureField: %+v\n", m)
			for _, sub := range m.Modules {
				switch t := sub.(type) {
				case *IOModuleConfigFill:
					fmt.Printf("IOSubModuleFill: %+v\n", t)
				case *IOModuleConfigFixed:
					fmt.Printf("IOSubModuleFixed: %+v\n", t)
				case *IOModuleConfigCalc:
					fmt.Printf("IOSubModuleCalc: %+v\n", t)
				case *IOModuleConfigCustom:
					fmt.Printf("IOSubModuleCustom: %+v\n", t)
				}
			}
		case *PrintActionFeatureField:
			fmt.Printf("PrintModuleFeatureField: %+v\n", m)
		case *DelayActionFeatureField:
			fmt.Printf("DelayModuleFeatureField: %+v\n", m)
		case *DeclareActionFeatureField:
			fmt.Printf("DeclareModuleFeatureField: %+v\n", m)
		case *IfActionFeatureField:
			fmt.Printf("IfModuleFeatureField: %+v\n", m)
		case *ElseActionFeatureField:
			fmt.Printf("ElseModuleFeatureField: %+v\n", m)
		case *ForActionFeatureField:
			fmt.Printf("ForModuleFeatureField: %+v\n", m)
		case *BlockEndActionFeatureField:
			fmt.Printf("BlockEndModuleFeatureField: %+v\n", m)
		case *LabelActionFeatureField:
			fmt.Printf("LabelModuleFeatureField: %+v\n", m)
		case *GotoActionFeatureField:
			fmt.Printf("GotoModuleFeatureField: %+v\n", m)
		case *ChangeBaudRateActionFeatureField:
			fmt.Printf("ChangeBaudRateModuleFeatureField: %+v\n", m)
		case *StopActionFeatureField:
			fmt.Printf("StopModuleFeatureField: %+v\n", m)

		default:
			println("no this type")
		}
	}
}

func ParseAppConfig(b []byte) (*AppConfig, error) {
	var appConfig AppConfig
	if err := json.Unmarshal(b, &appConfig); err != nil {
		return nil, err
	}
	return &appConfig, nil
}
