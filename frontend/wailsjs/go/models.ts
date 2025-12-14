export namespace main {
	
	export class AppSettings {
	    autoStart: boolean;
	    nginxPort: number;
	    phpPort: number;
	    mysqlPort: number;
	
	    static createFrom(source: any = {}) {
	        return new AppSettings(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.autoStart = source["autoStart"];
	        this.nginxPort = source["nginxPort"];
	        this.phpPort = source["phpPort"];
	        this.mysqlPort = source["mysqlPort"];
	    }
	}

}

export namespace services {
	
	export class RewriteRule {
	    id: number;
	    pattern: string;
	    destination: string;
	    type: string;
	    enabled: boolean;
	
	    static createFrom(source: any = {}) {
	        return new RewriteRule(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.pattern = source["pattern"];
	        this.destination = source["destination"];
	        this.type = source["type"];
	        this.enabled = source["enabled"];
	    }
	}
	export class ServiceInfo {
	    name: string;
	    version: string;
	    status: string;
	    port: string;
	
	    static createFrom(source: any = {}) {
	        return new ServiceInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.version = source["version"];
	        this.status = source["status"];
	        this.port = source["port"];
	    }
	}
	export class SiteConfig {
	    id: number;
	    name: string;
	    domain: string;
	    port: number;
	    root: string;
	    enabled: boolean;
	    rewriteRules: RewriteRule[];
	
	    static createFrom(source: any = {}) {
	        return new SiteConfig(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.domain = source["domain"];
	        this.port = source["port"];
	        this.root = source["root"];
	        this.enabled = source["enabled"];
	        this.rewriteRules = this.convertValues(source["rewriteRules"], RewriteRule);
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

}

