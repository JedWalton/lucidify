import { LocalStorage } from "./storage";

export interface ChangeLog {
    changeId?: number;
    key: keyof LocalStorage;
    operation: 'INSERT' | 'UPDATE' | 'DELETE';
    oldValue?: LocalStorage[keyof LocalStorage];
    newValue?: LocalStorage[keyof LocalStorage];
    timestamp: number;
}

