package bsd_testtool

type Config struct {
	AppName      string `json:"AppName"`
	SerialConfig struct {
		BaudRate int    `json:"BaudRate"`
		DataBits int    `json:"DataBits"`
		Parity   string `json:"Parity"`
		StopBits int    `json:"StopBits"`
	} `json:"SerialConfig"`
	DebugEnable       bool   `json:"DebugEnable"`
	LogEnable         bool   `json:"LogEnable"`
	LogExportEnable   bool   `json:"LogExportEnable"`
	LogExportLoaction string `json:"LogExportLoaction"`
	Actions           []struct {
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
		// 模块独有属性
		TypeFeatureField any `json:"TypeFeatureField"`
	}
}
