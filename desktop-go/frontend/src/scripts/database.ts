const DB_OBJECT: string = 'YekongaSystemDatabase';
const DB_NAME: string = 'DefaultSystemDatabase';

interface TableOptions {
    name: string;
    key?: string;
    indexes?: { name: string; key: string }[];
    relations?: { 
        [key: string]: { 
            table: string; 
            primaryKey: string;
            foreignKey: string;
            hasMany: boolean;
        }
    };
}

interface DatabaseOptions {
    version?: number;
    tables?: TableOptions[];
}

interface PaginationResult<T> {
    data: T[];
    total?: number;
    page?: number;
    pageSize?: number;
    hasMore: boolean;
}
 
export class YekongaDatabase {
    private request: IDBOpenDBRequest | undefined;
    private db_object: string = DB_OBJECT;
    private db_name: string;
    public options: DatabaseOptions;
    private db_version: number = 1;
    private tables: TableOptions[] | undefined;
    public db: IDBDatabase | undefined;

    constructor(name?: string | DatabaseOptions, options?: DatabaseOptions) {
        let _name: string | null = null;
        let _options: DatabaseOptions = {};

        if (typeof name === 'string') {
            _name = name;
            _options = options || {};
        } else {
            _name = null;
            _options = name || {};
        }

        this.db_name = _name || DB_NAME;
        this.options = _options;

        this.init(this.options);
    }

    public init(options: DatabaseOptions): void {
        if (typeof window.indexedDB === 'undefined') {
            window.indexedDB = (window as any).mozIndexedDB || (window as any).webkitIndexedDB || (window as any).msIndexedDB;
        }

        this.db_version = options?.version || 1;
        this.tables = Array.isArray(options?.tables) ? options.tables : [];

        if (window.indexedDB) {
            this.request = window.indexedDB.open(this.db_name, this.db_version);
            this.request.onupgradeneeded = this.onupgradeneeded();
            this.request.onsuccess = this.onsuccess();
            this.request.onerror = this.onerror();
        }
    }

    private onupgradeneeded(): (e: IDBVersionChangeEvent) => void {
        const $this = this;
        return (e: IDBVersionChangeEvent) => {
            const db = (e.target as IDBOpenDBRequest).result;
            if (Array.isArray($this.tables)) {
                for (const table of $this.tables) {
                    const test = true;
                    if (!db.objectStoreNames.contains(table.name) || test) {
                        const primaryKey = table.key || 'uuid';
                        const store = db.createObjectStore(table.name, { keyPath: primaryKey });

                        if (Array.isArray(table.indexes)) {
                            for (const ind of table.indexes) {
                                store.createIndex(ind.name, ind.key, { unique: false });
                            }
                        }
                    }
                }
            }
        };
    }

    private onsuccess(): (e: Event) => void {
        const $this = this;
        return (e: Event) => {
            $this.db = (e.target as IDBOpenDBRequest).result;
        };
    }

    private onerror(): (e: Event) => void {
        return (e: Event) => {
            console.error(`There was an error: ${(e.target as IDBRequest).error}`);
            console.error(e);
        };
    }

    public table<T>(name: string): Query<T> {
        return new Query<T>(name, this);
    }
}

export class Query<T>{
    private table: string;
    private parent: YekongaDatabase;
    private filters: { key: string; value: any; symbol: string }[] = [];
    private nextCursor?: IDBValidKey;
    private pageSize: number = 10;
    private orderByField: string = 'createdAt';
    private orderByDirection: 'asc' | 'desc' = 'asc';
    private relationKeys: string[] = [];

    readonly SYMBOL_EQUAL: string = '=';
    readonly SYMBOL_GREATER: string = '>';
    readonly SYMBOL_GREATER_EQUAL: string = '>=';
    readonly SYMBOL_SMALLER: string = '<';
    readonly SYMBOL_SMALLER_EQUAL: string = '<=';
    readonly SYMBOL_LIKE: string = 'LIKE';

    constructor(name: string, parent: YekongaDatabase) {
        this.table = name;
        this.parent = parent;
    }

    private db(): Promise<IDBDatabase | null> {
        const $parent = this.parent;
        $parent.init($parent.options);
        return new Promise((resolve) => {
            let i = 0;
            const timer = setInterval(() => {
                if ($parent.db) {
                    resolve($parent.db);
                    clearInterval(timer);
                }
                if (i > 100) {
                    resolve(null);
                    clearInterval(timer);
                }
                i++;
            }, 500);
        });
    }

    public async get(): Promise<T[]> {
        const $this = this;
        const db = await $this.db();
        let result: T[] = [];

        try {
            const tx = db!.transaction($this.table, 'readonly');
            const store = tx.objectStore($this.table);
            const request = store.getAll();

            tx.oncomplete = () => db!.close();

            result = await new Promise((resolve) =>
                request.onsuccess = (event) => {
                    return resolve($this.dataFilter((event.target as IDBRequest).result))
                }
            );
            
            const count = result.length;
            for (let i = 0; i < count; i++) {
                result[i] = await $this.processRelation(result[i]);
            }
        } catch (error) {
            result = [];
        }

        return result;
    }

    public async row(id: string | null = null): Promise<T | null> {
        const $this = this;
        const db = await $this.db();
        let result: T | null = null;

        try {
            const tx = db!.transaction($this.table, 'readonly');
            const store = tx.objectStore($this.table);
            tx.oncomplete = () => db!.close();
            result = await new Promise((resolve) => {
                if (id) {
                    store.get(id).onsuccess = (event) => resolve((event.target as IDBRequest).result);
                } else {
                    store.getAll().onsuccess = (event) => {
                        const res = $this.dataFilter((event.target as IDBRequest).result).shift();
                        resolve(res);
                    };
                }
            });

            if (result) result = await $this.processRelation(result);
        } catch (error) {
            result = null;
        }

        return result || null;
    }

    public async find(): Promise<T[]> {
        return await this.get();
    }

    public async findOne(id?: string | null): Promise<T | null> {
        return await this.row(id);
    }

    public async paginate(
        indexName?: string,
        pageSize?: number,
    ): Promise<PaginationResult<T>> {
        const $this = this;
        const db = await $this.db();
        $this.pageSize = pageSize ?? $this.pageSize;
        $this.orderByField = indexName ?? $this.orderByField;

        return new Promise((resolve, reject) => {
            const transaction = db!.transaction($this.table, 'readonly');
            const store = transaction.objectStore($this.table);
            const index = store.index($this.orderByField);

            const range = $this.nextCursor ? IDBKeyRange.lowerBound($this.nextCursor, true) : null;
            const request = index.openCursor(range);
            const data: T[] = [];

            request.onsuccess = async (event) => {
                const cur = (event.target as IDBRequest<IDBCursorWithValue>).result;

                if (cur && data.length < $this.pageSize) {
                    data.push(cur.value as T);
                    cur.continue();
                } else {
                    $this.nextCursor = cur ? cur.key : undefined;
                    // data.forEach(async (item, index) => {
                    //     data[index] = await $this.processRelation(item);
                    // });

                    const count = data.length;
                    for (let i = 0; i < count; i++) {
                        data[i] = await $this.processRelation(data[i]);
                    }

                    resolve({
                        data,
                        hasMore: !!cur,
                    });
                }
            };

            request.onerror = () => reject(request.error);
        });
    }

    private processRelation(value: any): T {
        for (const key of this.relationKeys) {
            const relation = this.parent.options.tables?.find(t => t.name === key)?.relations?.[key];
            if(relation) {
                const relatedTable = this.parent.table(relation.table);

                if(relation.hasMany) {
                    // primaryKey = "userId"; foreignKey = "id"
                    relatedTable.where(relation.primaryKey, value[relation.foreignKey]).get().then(res => {
                        (value as any)[key] = res;
                    });
                } else {
                    // primaryKey = "id"; foreignKey = "userId"
                    relatedTable.where(relation.primaryKey, value[relation.foreignKey]).row().then(res => {
                        (value as any)[key] = res;
                    });
                }
            }
        }

        return value;
    }

    public relations(...keys: string[]): this {
        this.relationKeys = keys;
        return this;
    }

    public ordersBy(field: string, direction: 'asc' | 'desc' = 'asc'): this {
        this.orderByField = field;
        this.orderByDirection = direction;

        return this;
    }

    public async create(data: any): Promise<{ status: string }> {
        const $this = this;
        const db = await $this.db();
        if (!data.uuid) data.uuid = $this.uuid();

        try {
            const tx = db!.transaction($this.table, 'readwrite');
            const store = tx.objectStore($this.table);

            store.put(data);
            tx.oncomplete = () => db!.close();
        } catch (error) {
            console.error(error);
            return { status: 'fail' };
        }

        return { status: 'success' };
    }

    public async update(data: any): Promise<{ status: string }> {
        return this.create(data);
    }

    public async delete(id: string): Promise<{ status: string }> {
        const db = await this.db();
        try {
            const tx = db!.transaction(this.table, 'readwrite');
            const store = tx.objectStore(this.table);

            tx.oncomplete = () => db!.close();
            store.delete(id);
        } catch (error) {
            console.error(error);
        }

        return { status: 'success' };
    }

    public whereLike(key: string, arg: any = null): this {
        return this.where(key, 'LIKE', arg);
    }

    public where(key: string, arg1?: string | any, arg2?: any): this {
        let value: any = null;
        let symbol: string = '=';

        if (arg2 === undefined) {
            value = arg1;
        } else {
            value = arg2;
            symbol = arg1;
        }

        if (!Array.isArray(this.filters)) this.filters = [];

        this.filters.push({
            key,
            value,
            symbol,
        });

        return this;
    }

    public async clear(): Promise<{ status: string }> {
        try {
            const db = await this.db();
            const tx = db!.transaction(this.table, 'readwrite');
            const store = tx.objectStore(this.table);

            tx.oncomplete = () => db!.close();
            store.clear();
        } catch (error) { }

        return { status: 'success' };
    }

    private dataFilter(result: any[]): any[] {
        const filters: { [key: string]: { key: string; value: any; symbol: string }[] } = {};
        if (Array.isArray(this.filters)) {
            for (const item of this.filters) {
                if (!Array.isArray(filters[item.key])) filters[item.key] = [];
                filters[item.key].push(item);
            }
        }

        return result.filter((elem) => {
            for (const key in elem) {
                if (Object.prototype.hasOwnProperty.call(elem, key) && Object.prototype.hasOwnProperty.call(filters, key)) {
                    const comp = filters[key][0];
                    let v1 = elem[key];
                    let v2 = comp.value;
                    const symbol = comp.symbol;

                    if (typeof v1 === 'string') {
                        v1 = v1.toLowerCase();
                        v2 = v2.toLowerCase();
                    }

                    if (symbol === this.SYMBOL_EQUAL) {
                        if (!(v1 === v2)) return false;
                    } else if (symbol === this.SYMBOL_GREATER) {
                        if (!(v1 > v2)) return false;
                    } else if (symbol === this.SYMBOL_GREATER_EQUAL) {
                        if (!(v1 >= v2)) return false;
                    } else if (symbol === this.SYMBOL_SMALLER) {
                        if (!(v1 < v2)) return false;
                    } else if (symbol === this.SYMBOL_SMALLER_EQUAL) {
                        if (!(v1 <= v2)) return false;
                    } else if (symbol === this.SYMBOL_LIKE) {
                        if (!v1.includes(v2)) return false;
                    }
                }
            }
            return true;
        });
    }

    private uuid(): string {
        const dt = new Date().getTime();
        return 'xxxxxxxx-xxxx-4xxx-yxxx-xxxxxxxxxxxx'.replace(/[xy]/g, (c) => {
            const r = (dt + Math.random() * 16) % 16 | 0;
            return (c === 'x' ? r : (r & 0x3) | 0x8).toString(16);
        });
    }
}

window.VueDatabase = (name?: string | DatabaseOptions, options?: DatabaseOptions) => new YekongaDatabase(name, options)
window.YekongaSystemDatabase = (name?: string | DatabaseOptions, options?: DatabaseOptions) => new YekongaDatabase(name, options)

declare global {
    interface Window {
        VueDatabase: (name?: string | DatabaseOptions, options?: DatabaseOptions) => YekongaDatabase;
        YekongaSystemDatabase: (name?: string | DatabaseOptions, options?: DatabaseOptions) => YekongaDatabase;
    }
}

export default YekongaDatabase;