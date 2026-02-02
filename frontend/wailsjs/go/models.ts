export namespace models {
	
	export class AppConfig {
	    autoRefresh: boolean;
	    refreshInterval: number;
	    activeGroups: string[];
	    backupEnabled: boolean;
	    maxBackups: number;
	    systemHostPath: string;
	    createdAt: string;
	    updatedAt: string;
	
	    static createFrom(source: any = {}) {
	        return new AppConfig(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.autoRefresh = source["autoRefresh"];
	        this.refreshInterval = source["refreshInterval"];
	        this.activeGroups = source["activeGroups"];
	        this.backupEnabled = source["backupEnabled"];
	        this.maxBackups = source["maxBackups"];
	        this.systemHostPath = source["systemHostPath"];
	        this.createdAt = source["createdAt"];
	        this.updatedAt = source["updatedAt"];
	    }
	}
	export class HostGroup {
	    id: string;
	    name: string;
	    description?: string;
	    content: string;
	    enabled: boolean;
	    isRemote: boolean;
	    url?: string;
	    refreshInterval: number;
	    lastUpdated: string;
	    createdAt: string;
	    updatedAt: string;
	
	    static createFrom(source: any = {}) {
	        return new HostGroup(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.description = source["description"];
	        this.content = source["content"];
	        this.enabled = source["enabled"];
	        this.isRemote = source["isRemote"];
	        this.url = source["url"];
	        this.refreshInterval = source["refreshInterval"];
	        this.lastUpdated = source["lastUpdated"];
	        this.createdAt = source["createdAt"];
	        this.updatedAt = source["updatedAt"];
	    }
	}

}

