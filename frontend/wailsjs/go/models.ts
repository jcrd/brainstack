export namespace main {
	
	export class Task {
	    stack_id: number;
	    text: string;
	    done: boolean;
	    order: number;
	
	    static createFrom(source: any = {}) {
	        return new Task(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.stack_id = source["stack_id"];
	        this.text = source["text"];
	        this.done = source["done"];
	        this.order = source["order"];
	    }
	}
	export class Stack {
	    name: string;
	    order: number;
	    tasks: Task[];
	
	    static createFrom(source: any = {}) {
	        return new Stack(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.order = source["order"];
	        this.tasks = this.convertValues(source["tasks"], Task);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice) {
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

