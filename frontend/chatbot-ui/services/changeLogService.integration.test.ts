// changeLogService.test.ts
import { changeLogService } from './changeLogService';
import { ChangeLog } from '@/types/changelog';
import { describe, expect, it, beforeEach } from 'vitest';

global.localStorage = (() => {
    let store: { [key: string]: string } = {};
    return {
        getItem: (key: string) => store[key] || null,
        setItem: (key: string, value: string) => { store[key] = value },
        removeItem: (key: string) => { delete store[key] },
        clear: () => { store = {} },
    };
})();


// Only mock if window is defined
if (typeof window !== 'undefined') {
    Object.defineProperty(window, 'localStorage', {
        value: localStorage
    });
}


describe('changeLogService', () => {
    beforeEach(() => {
        localStorage.clear();
    });

    it('should get empty changelog if none exists', () => {
        const logs = changeLogService.getChangeLog();
        expect(logs).toEqual([]);
    });

    it('should add to changelog and generate new changeId', () => {
        const change: Omit<ChangeLog, 'changeId'> = {
            key: 'apiKey',
            operation: 'INSERT',
            newValue: 'new-api-key',
            timestamp: Date.now()
        };

        changeLogService.addToChangeLog(change);

        const logs = changeLogService.getChangeLog();
        expect(logs).toHaveLength(1);
        expect(logs[0].changeId).toBe(1);
        expect(logs[0].newValue).toBe('new-api-key');
    });

    it('should remove item from changelog by changeId', () => {
        const change: Omit<ChangeLog, 'changeId'> = {
            key: 'apiKey',
            operation: 'INSERT',
            newValue: 'new-api-key',
            timestamp: Date.now()
        };

        changeLogService.addToChangeLog(change);
        changeLogService.removeFromChangeLog(1);

        const logs = changeLogService.getChangeLog();
        expect(logs).toEqual([]);
    });

    it('should clear the changelog', () => {
        const change: Omit<ChangeLog, 'changeId'> = {
            key: 'apiKey',
            operation: 'INSERT',
            newValue: 'new-api-key',
            timestamp: Date.now()
        };

        changeLogService.addToChangeLog(change);
        changeLogService.clearChangeLog();

        const logs = changeLogService.getChangeLog();
        expect(logs).toEqual([]);
    });
});

