// storageService.integration.test.ts
import { storageService } from '@/services/storageService';
import { fail, describe, expect, it, beforeEach } from 'vitest';
import { LocalStorage } from '@/types/storage';

// // Mock local storage
// let store: any = {};
//
// beforeEach(() => {
//   store = {};
//
//   let localStorageMock = {
//     getItem: function(key: string) {
//       return store[key] || null;
//     },
//     setItem: function(key: string, value: string) {
//       store[key] = value
//     },
//     removeItem: function(key: string) {
//       delete store[key];
//     }
//   };
//
//   if (typeof window !== 'undefined') {
//     // If window is available, set the mock on it
//     Object.defineProperty(window, 'localStorage', {
//       value: localStorageMock
//     });
//   } else {
//     // If window is not available (e.g., running in Node.js), use global
//     (global as any).localStorage = localStorageMock;
//   }
// });

describe('storageService set and get', () => {
  const testKey = 'apiKey' as keyof LocalStorage;
  // const testValue = 'testValue' as LocalStorage[keyof LocalStorage];
  process.env.PUBLIC_BACKEND_API_URL = 'http://localhost:8080';

  it('should set and get an item', async () => {
    let resp = await storageService.setItemWrapper(testKey, 'testValue');

    if (resp === null) {
      fail('resp is null');
    }
    if (resp !== null) {
      const parsedResp: any = JSON.parse(resp);
      expect(parsedResp.success).toBe(true);
    }

    let respGet = await storageService.getItemWrapper(testKey)
    if (respGet === null) {
      fail('respGet is null');
    }
    if (respGet !== null) {
      const parsedRespGet: any = JSON.parse(respGet);
      expect(parsedRespGet.success).toBe(true);
    }

  });

  // it('should set item, remove and verify removed through get', async () => {
  //   await storageService.setItem(testKey, 'testValue');
  //   let item = await storageService.getItem(testKey);
  // });
});

