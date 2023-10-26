// storageService.integration.test.ts
import { storageService } from '@/services/storageService';
import { describe, expect, it, beforeEach } from 'vitest';
import { LocalStorage } from '@/types/storage';

// Mock local storage
let store: any = {};

beforeEach(() => {
  store = {};

  let localStorageMock = {
    getItem: function (key: string) {
      return store[key] || null;
    },
    setItem: function (key: string, value: string) {
      store[key] = value.toString();
    },
    removeItem: function (key: string) {
      delete store[key];
    }
  };

  if (typeof window !== 'undefined') {
    // If window is available, set the mock on it
    Object.defineProperty(window, 'localStorage', {
      value: localStorageMock
    });
  } else {
    // If window is not available (e.g., running in Node.js), use global
    (global as any).localStorage = localStorageMock;
  }
});

describe('storageService - changeLog functions', () => {
  it('should return empty change log if none exists', () => {
    const log = storageService.getChangeLog();
    expect(log).toEqual([]);
  });

  it('should return correct change log if it exists', () => {
    const dummyLog = [{ changeId: 1, key: 'testKey', newValue: 'testValue', operation: 'INSERT' }];
    localStorage.setItem(CHANGE_LOG_KEY, JSON.stringify(dummyLog));
    const log = getChangeLog();
    expect(log).toEqual(dummyLog);
  });

  it('should add an entry to an empty change log', () => {
    addToChangeLog({ key: 'testKey', newValue: 'testValue', operation: 'INSERT' });
    const log = getChangeLog();
    expect(log).toEqual([{ changeId: 1, key: 'testKey', newValue: 'testValue', operation: 'INSERT' }]);
  });

  it('should have incrementing change IDs for consecutive entries', () => {
    addToChangeLog({ key: 'testKey2', newValue: 'testValue2', operation: 'UPDATE' });
    const log = getChangeLog();
    expect(log[1].changeId).toBe(2);
  });

  it('should clear the change log', () => {
    clearChangeLog();
    const log = getChangeLog();
    expect(log).toEqual([]);
  });

  it('should remove specified entry from the change log', () => {
    addToChangeLog({ key: 'testKey', newValue: 'testValue', operation: 'INSERT' });
    const initialLog = getChangeLog();
    removeFromChangeLog(initialLog[0].changeId);
    const updatedLog = getChangeLog();
    expect(updatedLog).toEqual([]);
  });

  it('should do nothing if trying to remove a nonexistent entry', () => {
    addToChangeLog({ key: 'testKey', newValue: 'testValue', operation: 'INSERT' });
    removeFromChangeLog(9999); // Using a nonexistent changeId
    const log = getChangeLog();
    expect(log).toHaveLength(1);
  });
});

describe('storageService', () => {
  const testKey = 'testKey' as keyof LocalStorage;
  // const testValue = 'testValue' as LocalStorage[keyof LocalStorage];

  it('should add and retrieve an item', async () => {
    await storageService.setItem(testKey, 'testValue');
    expect(await storageService.getItem(testKey)).toBe('testValue');
  });

  it('should update change log on setItem', async () => {
    await storageService.setItem(testKey, 'testValue');
    const changeLog = JSON.parse(store['__CHANGE_LOG__']);
    expect(changeLog[0].operation).toBe('INSERT');
    expect(changeLog[0].key).toBe('testKey');
    expect(changeLog[0].newValue).toBe('testValue');

    await storageService.setItem(testKey, 'newValue');
    const updatedChangeLog = JSON.parse(store['__CHANGE_LOG__']);
    expect(updatedChangeLog[1].operation).toBe('UPDATE');
    expect(updatedChangeLog[1].key).toBe('testKey');
    expect(updatedChangeLog[1].newValue).toBe('newValue');
  });

  it('should remove an item and update change log', async () => {
    await storageService.setItem(testKey, 'testValue');
    await storageService.removeItem(testKey);
    expect(await storageService.getItem(testKey)).toBeNull();

    const changeLog = JSON.parse(store['__CHANGE_LOG__']);
    expect(changeLog[1].operation).toBe('DELETE');
    expect(changeLog[1].key).toBe('testKey');
  });

});

