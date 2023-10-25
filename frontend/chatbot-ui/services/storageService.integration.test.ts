// storageService.integration.test.ts
import { storageService } from '@/services/storageService';
import { describe, expect, it, afterAll } from 'vitest';

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
      store[key] = value.toString(); // toString() is redundant here since value is already a string
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

  it('syncs individual changes with the server', async () => {
    await storageService.setItem('apiKey', 'valueToSync');
    const changeLogBefore = JSON.parse(localStorage.getItem('__CHANGE_LOG__')!);
    await storageService.syncSingleChangeWithServer(changeLogBefore[changeLogBefore.length - 1]);
    const changeLogAfter = JSON.parse(localStorage.getItem('__CHANGE_LOG__')!);
    expect(changeLogAfter.length).toBe(changeLogBefore.length - 1);
  });

  it('syncs all changes with the server', async () => {
    await storageService.setItem('apiKey', 'firstValue');
    await storageService.setItem('folders', 'secondValue');
    const changeLogBefore = JSON.parse(localStorage.getItem('__CHANGE_LOG__')!);
    expect(changeLogBefore.length).toBeGreaterThan(1);
    await storageService.syncAllChangesWithServer();
    const changeLogAfter = JSON.parse(localStorage.getItem('__CHANGE_LOG__')!);
    expect(changeLogAfter.length).toBe(0);
  });

});

describe('storageService Integration Tests', () => {
  const testKey = 'apiKey';
  const testValue = 'integration_test_value';

  afterAll(async () => {
    // Clean up the test data from the server and local storage
    await storageService.removeItem(testKey);
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
    localStorage.removeItem(testKey); // Ensure the item is not in local storage

    const value = await storageService.getItem(testKey);
    expect(value).toBe(`{"success":true,"data":{"exampleKey":"exampleValue"},"message":"Data fetched successfully"}`); // Server should return the original value
  });

  it('removeItem - removes item from local storage and server', async () => {
    await storageService.setItem(testKey, testValue); // Ensure the item is set before removing it
    await storageService.removeItem(testKey);

    expect(localStorage.getItem(testKey)).toBeNull();

    // Verify that the item is also removed from the server
    const valueFromServer = await storageService.getItem(testKey);
    expect(valueFromServer).toBeNull();
  });

  // ...additional tests...
});

