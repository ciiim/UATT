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
	export class AppConfigSettings {
	    AppName: string;
	    // Go type: struct { BaudRate int "json:\"BaudRate\""; DataBits int "json:\"DataBits\""; Parity string "json:\"Parity\""; StopBits int "json:\"StopBits\"" }
	    SerialConfig: any;
	    LogEnable: boolean;
	    LogExportEnable: boolean;
	    LogExportLoaction: string;
	
	    static createFrom(source: any = {}) {
	        return new AppConfigSettings(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.AppName = source["AppName"];
	        this.SerialConfig = this.convertValues(source["SerialConfig"], Object);
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

