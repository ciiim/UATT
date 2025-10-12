export interface ConfigActionBase {
  ActionUID: number
  ActionType: string
  ActionTypeID: number
  Name: string
  BreakPoint: boolean
  TypeFeatureField: any
  Tags: any[]
  Status: string
}

export interface Tag {
  label: string;
  len: number;
}

export interface SerialConfig {
  BaudRate: number;
	DataBits: number;
	Parity: string;
	StopBits: number;
}

export interface AppConfigSettings {
  AppName: string;
	SerialConfig: SerialConfig;
	LogEnable: boolean;
	LogExportEnable: boolean;
	LogExportLoaction: string;
}

export interface ActionReport {
  ActionName : string;
  ActionUID : number;
  Result : string;
}