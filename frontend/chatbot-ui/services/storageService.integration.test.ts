// storageService.integration.test.ts
import { storageService } from '@/services/storageService';
import { describe, expect, it, afterAll, beforeEach } from 'vitest';

// Define an interface for the mock storage, which replicates the Storage interface
interface IMockStorage {
  [key: string]: string;
}

const localStorageMock = (function() {
  let store: IMockStorage = {};
  return {
    getItem: function(key: string): string | null {
      return store[key] || null;
    },
    setItem: function(key: string, value: string): void {
      store[key] = value;
    },
    removeItem: function(key: string): void {
      delete store[key];
    },
    clear: function(): void {
      store = {};
    },
  };
})();

// This casting is necessary to allow TypeScript to let you add the localStorage to the global object
(global as any).localStorage = localStorageMock;

describe('storageService Integration Tests - Server Sync', () => {

  beforeEach(() => {
    // Clear local storage before each test for a consistent starting point
    localStorage.clear();
  });

  const testKey = 'apiKey';
  const testValue = 'exampleValue';

  it('setItem - sets item in local storage and syncs with server', async () => {
    await storageService.setItem(testKey, testValue);

    const locallyStoredValue = localStorage.getItem(testKey);
    expect(locallyStoredValue).toBe(testValue); // Assuming you're storing JSON strings

    const valueFromServer = await storageService.getItem(testKey);
    expect(valueFromServer).toBe(testValue);
  });

  it('getItem - retrieves item from server when not in local storage', async () => {
    expect(localStorage.getItem(testKey)).toBeNull();
    await storageService.setItem(testKey, testValue);
    const value = await storageService.getItem(testKey);
    if (!value) {
      throw new Error('Expected value to not be null.');
    }

    expect(value).toBe("exampleValue");
    const locallyStoredValue = localStorage.getItem(testKey);
    expect(locallyStoredValue).toBe("exampleValue");
  });

  it('syncs individual changes with the server', async () => {
    await storageService.setItem('apiKey', 'valueToSync');
    const changeLogBefore = JSON.parse(localStorage.getItem('__CHANGE_LOG__')!);
    expect(changeLogBefore).toBeDefined(); // Add this check
    expect(changeLogBefore.length).toBeGreaterThan(0); // Check that there's at least one entry
    await storageService.syncSingleChangeWithServer(changeLogBefore[changeLogBefore.length - 1]);
    const changeLogAfter = JSON.parse(localStorage.getItem('__CHANGE_LOG__')!);
    expect(changeLogAfter.length).toBe(changeLogBefore.length);
  });

  it('syncs all changes with the server', async () => {
    await storageService.setItem('apiKey', 'firstValue');
    await storageService.setItem('folders', 'secondValue');
    const changeLogBefore = JSON.parse(localStorage.getItem('__CHANGE_LOG__')!);
    expect(changeLogBefore.length).toBe(2);
    await storageService.syncAllChangesWithServer();
    const changeLogAfter = JSON.parse(localStorage.getItem('__CHANGE_LOG__')!);
    expect(changeLogAfter.length).toBe(0);
  });

});

describe('storageService Integration Tests', () => {
  const testKey = 'apiKey';
  const testValue = 'storageService integration test value';

  beforeEach(async () => {
    // Clean up the test data from the server and local storage
    await storageService.removeItem(testKey);
    localStorage.removeItem(testKey);
  });

  it('setItem - sets item in local storage and syncs with server', async () => {
    await storageService.setItem(testKey, testValue);
    expect(localStorage.getItem(testKey)).toBe(testValue);

    // Verify that the item is also set on the server
    // Here you might need a way to directly query your server's data store, which depends on your implementation
    const valueFromServer = await storageService.getItem(testKey);
    expect(valueFromServer).toBe(testValue); // Adjust based on how your server responds
  });

  it('getItem - retrieves item from server when not in local storage', async () => {
    await storageService.setItem(testKey, testValue);
    expect(localStorage.getItem(testKey)).toBe(testValue);
    localStorage.removeItem(testKey); // Ensure the item is not in local storage
    expect(localStorage.getItem(testKey)).toBe(null);

    const value = await storageService.getItem(testKey);
    if (value === null) {
      // Handle the null case, maybe throw an error or provide a default value
      throw new Error("Value not found in localStorage");
    }
    // storageService.syncAllChangesWithServer()
    expect(value).toBe(testValue);

    expect(localStorage.getItem(testKey)).toBe(testValue);
  });

  it('removeItem - removes item from local storage and this should be reflected on server', async () => {
    await storageService.setItem(testKey, testValue); // Ensure the item is set before removing it
    expect(localStorage.getItem(testKey)).toBe(testValue);
    await storageService.removeItem(testKey);
    expect(localStorage.getItem(testKey)).toBeNull();
    // Verify that the item is also removed from the server
    const valueFromServer = await storageService.getItem(testKey);
    expect(valueFromServer).toBeNull(); // Adjust based on how your server responds
  });

  // ...additional tests...
});

