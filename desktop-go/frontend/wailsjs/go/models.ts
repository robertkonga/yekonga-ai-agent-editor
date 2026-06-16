export namespace types {
	
	export class FileNode {
	    id: string;
	    name: string;
	    path: string;
	    extension?: string;
	    type: string;
	    lang?: string;
	    expanded: boolean;
	    children?: FileNode[];
	    // Go type: time
	    lastUpdate: any;
	
	    static createFrom(source: any = {}) {
	        return new FileNode(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.path = source["path"];
	        this.extension = source["extension"];
	        this.type = source["type"];
	        this.lang = source["lang"];
	        this.expanded = source["expanded"];
	        this.children = this.convertValues(source["children"], FileNode);
	        this.lastUpdate = this.convertValues(source["lastUpdate"], null);
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

