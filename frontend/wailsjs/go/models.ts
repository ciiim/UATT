export namespace bsd_testtool {
	
	export class App {
	    AppName: string;
	
	    static createFrom(source: any = {}) {
	        return new App(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.AppName = source["AppName"];
	    }
	}
	export class SerialConfig {
	    BaudRate: number;
	    DataBits: number;
	    Parity: string;
	    StopBits: number;
	
	    static createFrom(source: any = {}) {
	        return new SerialConfig(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.BaudRate = source["BaudRate"];
	        this.DataBits = source["DataBits"];
	        this.Parity = source["Parity"];
	        this.StopBits = source["StopBits"];
	    }
	}
	export class AppConfigSettings {
	    AppName: string;
	    SerialConfig: SerialConfig;
	    LogEnable: boolean;
	    LogExportEnable: boolean;
	    LogExportLoaction: string;
	
	    static createFrom(source: any = {}) {
	        return new AppConfigSettings(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.AppName = source["AppName"];
	        this.SerialConfig = this.convertValues(source["SerialConfig"], SerialConfig);
	        this.LogEnable = source["LogEnable"];
	        this.LogExportEnable = source["LogExportEnable"];
	        this.LogExportLoaction = source["LogExportLoaction"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class Position {
	    X: number;
	    Y: number;
	
	    static createFrom(source: any = {}) {
	        return new Position(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.X = source["X"];
	        this.Y = source["Y"];
	    }
	}
	export class CanvasComponent {
	    ID: string;
	    Type: string;
	    Label: string;
	    AttachApp: string;
	    Value: string;
	    Position: Position;
	
	    static createFrom(source: any = {}) {
	        return new CanvasComponent(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.ID = source["ID"];
	        this.Type = source["Type"];
	        this.Label = source["Label"];
	        this.AttachApp = source["AttachApp"];
	        this.Value = source["Value"];
	        this.Position = this.convertValues(source["Position"], Position);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class CanvasComponentConnection {
	    FromID: string;
	    ToID: string;
	
	    static createFrom(source: any = {}) {
	        return new CanvasComponentConnection(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.FromID = source["FromID"];
	        this.ToID = source["ToID"];
	    }
	}
	export class CanvasData {
	    ComponentList: CanvasComponent[];
	    Connections: CanvasComponentConnection[];
	
	    static createFrom(source: any = {}) {
	        return new CanvasData(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.ComponentList = this.convertValues(source["ComponentList"], CanvasComponent);
	        this.Connections = this.convertValues(source["Connections"], CanvasComponentConnection);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class CanvasConfig {
	    CanvasName: string;
	    Data: CanvasData;
	
	    static createFrom(source: any = {}) {
	        return new CanvasConfig(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.CanvasName = source["CanvasName"];
	        this.Data = this.convertValues(source["Data"], CanvasData);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	
	export class ConfigActionBaseJson {
	    ActionUID: number;
	    ActionType: string;
	    ActionTypeID: number;
	    Name: string;
	    BreakPoint: boolean;
	    TypeFeatureField: number[];
	
	    static createFrom(source: any = {}) {
	        return new ConfigActionBaseJson(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.ActionUID = source["ActionUID"];
	        this.ActionType = source["ActionType"];
	        this.ActionTypeID = source["ActionTypeID"];
	        this.Name = source["Name"];
	        this.BreakPoint = source["BreakPoint"];
	        this.TypeFeatureField = source["TypeFeatureField"];
	    }
	}
	

}

