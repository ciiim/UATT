package bsd_testtool

import (
	"encoding/json"
	"fmt"
)

type ModuleUID int

type App struct {

	// 可以只解析APPName
	AppName string `json:"AppName"`

	// 存放位置
	appFileLocation string

	config *AppConfig
}

type AppConfig struct {
	AppName      string `json:"AppName"`
	SerialConfig struct {
		BaudRate int    `json:"BaudRate"`
		DataBits int    `json:"DataBits"`
		Parity   string `json:"Parity"`
		StopBits int    `json:"StopBits"`
	} `json:"SerialConfig"`
	DebugEnable       bool         `json:"DebugEnable"`
	LogEnable         bool         `json:"LogEnable"`
	LogExportEnable   bool         `json:"LogExportEnable"`
	LogExportLoaction string       `json:"LogExportLoaction"`
	Actions           []ModuleBase `json:"Actions"`
}

func (a *App) PrintConfig() {
	fmt.Printf("app: %v\n", a.config)
	for _, mod := range a.config.Actions {
		switch m := mod.TypeFeatureField.(type) {
		case *IOModuleFeatureField:
			fmt.Printf("IOModuleFeatureField: %+v\n", m)
			for _, sub := range m.SubModules {
				switch t := sub.(type) {
				case *IOSubModuleFill:
					fmt.Printf("IOSubModuleFill: %+v\n", t)
				case *IOSubModuleFixed:
					fmt.Printf("IOSubModuleFixed: %+v\n", t)
				case *IOSubModuleCalc:
					fmt.Printf("IOSubModuleCalc: %+v\n", t)
				case *IOSubModuleCustom:
					fmt.Printf("IOSubModuleCustom: %+v\n", t)
				}
			}
		case *PrintModuleFeatureField:
			fmt.Printf("PrintModuleFeatureField: %+v\n", m)
		case *DelayModuleFeatureField:
			fmt.Printf("DelayModuleFeatureField: %+v\n", m)
		case *DeclareModuleFeatureField:
			fmt.Printf("DeclareModuleFeatureField: %+v\n", m)
		case *IfModuleFeatureField:
			fmt.Printf("IfModuleFeatureField: %+v\n", m)
		case *ElseModuleFeatureField:
			fmt.Printf("ElseModuleFeatureField: %+v\n", m)
		case *ForModuleFeatureField:
			fmt.Printf("ForModuleFeatureField: %+v\n", m)
		case *BlockEndModuleFeatureField:
			fmt.Printf("BlockEndModuleFeatureField: %+v\n", m)
		case *LabelModuleFeatureField:
			fmt.Printf("LabelModuleFeatureField: %+v\n", m)
		case *GotoModuleFeatureField:
			fmt.Printf("GotoModuleFeatureField: %+v\n", m)
		case *ChangeBaudRateModuleFeatureField:
			fmt.Printf("ChangeBaudRateModuleFeatureField: %+v\n", m)
		case *StopModuleFeatureField:
			fmt.Printf("StopModuleFeatureField: %+v\n", m)

		default:
			println("no this type")
		}
	}
}

func (a *App) StaticCheck() []error {
	return nil
}

func ParseAppConfig(b []byte) (*AppConfig, error) {
	var appConfig AppConfig
	if err := json.Unmarshal(b, &appConfig); err != nil {
		return nil, err
	}
	return &appConfig, nil
}
