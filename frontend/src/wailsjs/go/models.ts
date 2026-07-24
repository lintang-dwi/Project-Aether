export namespace ui {
	
	export class AuditDTO {
	    time: string;
	    perm: string;
	    target: string;
	    allowed: boolean;
	    reason?: string;
	
	    static createFrom(source: any = {}) {
	        return new AuditDTO(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.time = source["time"];
	        this.perm = source["perm"];
	        this.target = source["target"];
	        this.allowed = source["allowed"];
	        this.reason = source["reason"];
	    }
	}
	export class EdgeDTO {
	    id: string;
	    from: string;
	    to: string;
	    type: string;
	    weight: number;
	
	    static createFrom(source: any = {}) {
	        return new EdgeDTO(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.from = source["from"];
	        this.to = source["to"];
	        this.type = source["type"];
	        this.weight = source["weight"];
	    }
	}
	export class FileDTO {
	    path: string;
	    size: number;
	    modTime: string;
	
	    static createFrom(source: any = {}) {
	        return new FileDTO(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.path = source["path"];
	        this.size = source["size"];
	        this.modTime = source["modTime"];
	    }
	}
	export class NodeDTO {
	    id: string;
	    type: string;
	    name: string;
	    path: string;
	    language?: string;
	
	    static createFrom(source: any = {}) {
	        return new NodeDTO(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.type = source["type"];
	        this.name = source["name"];
	        this.path = source["path"];
	        this.language = source["language"];
	    }
	}
	export class GraphDataDTO {
	    nodes: NodeDTO[];
	    edges: EdgeDTO[];
	
	    static createFrom(source: any = {}) {
	        return new GraphDataDTO(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.nodes = this.convertValues(source["nodes"], NodeDTO);
	        this.edges = this.convertValues(source["edges"], EdgeDTO);
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
	
	export class PermissionStatusDTO {
	    allowWrite: boolean;
	    allowExecute: boolean;
	
	    static createFrom(source: any = {}) {
	        return new PermissionStatusDTO(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.allowWrite = source["allowWrite"];
	        this.allowExecute = source["allowExecute"];
	    }
	}
	export class TaskDTO {
	    id: string;
	    goal: string;
	    status: string;
	    targetPath: string;
	    createdAt: string;
	    error?: string;
	
	    static createFrom(source: any = {}) {
	        return new TaskDTO(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.goal = source["goal"];
	        this.status = source["status"];
	        this.targetPath = source["targetPath"];
	        this.createdAt = source["createdAt"];
	        this.error = source["error"];
	    }
	}

}

