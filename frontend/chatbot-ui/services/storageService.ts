// storageService.ts
import { LocalStorage } from '@/types/storage';
import { changeLogService } from './changeLogService';


export const storageService = {
  async getItem(key: keyof LocalStorage): Promise<string> {
    // localStorage.getItem(key) returns null if the key does not exist
    return await this.getItemFromServer(key);
  },

  async setItem(key: keyof LocalStorage, value: LocalStorage[keyof LocalStorage]): Promise<string | null> {
    //   localStorage.setItem(key, String(value));
    return await this.setItemOnServer(key, value);

  },

  async removeItem(key: keyof LocalStorage): Promise<string | null> {
    // localStorage.removeItem(key);
    return await this.removeItemFromServer(key);
  },

  async getItemFromServer(key: keyof LocalStorage): Promise<string> {
    try {
      const url = `${process.env.PUBLIC_BACKEND_API_URL}/api/sync/localstorage/?key=${encodeURIComponent(key as string)}`;

      const options: RequestInit = {
        method: 'GET',
        headers: {
          'Content-Type': 'application/json',
        },
        mode: 'cors',
      };

      const response = await fetch(url, options);

      if (!response.ok) {
        // You can first attempt to decode the response as JSON, and then fall back to text if it fails.
        let errorMessage = 'Server responded with an error';
        try {
          const errorBody = await response.json();
          errorMessage = errorBody.message || `Server responded with ${response.status}`;
        } catch (jsonError) {
          errorMessage = await response.text(); // If response is not in JSON format
        }

        throw new Error(errorMessage);
      }

      // If the response is OK, we decode it from JSON
      const data = await response.json();
      return JSON.stringify(data);  // or just `return data;` if you don't need to stringify the response
    } catch (error) {
      console.error(`Request failed: ${error}`);
      throw Error();
    }
  },

  async setItemOnServer(key: keyof LocalStorage, value: LocalStorage[keyof LocalStorage]): Promise<string | null> {
    try {
      const url = `${process.env.PUBLIC_BACKEND_API_URL}/api/sync/localstorage?key=${encodeURIComponent(key as string)}`;

      const options: RequestInit = {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        mode: 'cors',
        body: JSON.stringify({ value }), // We send only the value in the body as the key is already in the URL.
      };

      const response = await fetch(url, options);
      const responseClone = response.clone(); // Clone the response to read it multiple times

      if (!response.ok) {
        let errorMessage = 'Server responded with an error';
        try {
          const errorBody = await responseClone.json(); // Try to parse as JSON first
          errorMessage = errorBody.message || `Server responded with status code ${response.status}`;
        } catch (jsonError) {
          errorMessage = await responseClone.text(); // If response is not JSON, read as text
        }

        throw new Error(errorMessage);
      }

      // handle response logic, if needed
      const data = await response.json();
      return JSON.stringify(data);  // or just `return data;` if you don't need to stringify the response

    } catch (error) {
      console.error('Failed to sync with server:', error);
      return null;
    }
  },

  async removeItemFromServer(key: keyof LocalStorage): Promise<string | null> {
    try {
      const url = `${process.env.PUBLIC_BACKEND_API_URL}/api/sync/localstorage/?key=${encodeURIComponent(key as string)}`;

      const options: RequestInit = {
        method: 'GET',
        headers: {
          'Content-Type': 'application/json',
        },
        mode: 'cors',
      };

      const response = await fetch(url, options);

      if (!response.ok) {
        // You can first attempt to decode the response as JSON, and then fall back to text if it fails.
        let errorMessage = 'Server responded with an error';
        try {
          const errorBody = await response.json();
          errorMessage = errorBody.message || `Server responded with ${response.status}`;
        } catch (jsonError) {
          errorMessage = await response.text(); // If response is not in JSON format
        }

        throw new Error(errorMessage);
      }

      // If the response is OK, we decode it from JSON
      const data = await response.json();
      return JSON.stringify(data);  // or just `return data;` if you don't need to stringify the response
    } catch (error) {
      console.error(`Request failed: ${error}`);
      return null;
    }
  },
};
//   async syncAllChangesWithServer(): Promise<void> {
//     const changeLog = getChangeLog() as ChangeLog[];
//     if (!changeLog || !changeLog.length) {
//       throw new Error("ChangeLog is empty or not valid");
//     }
//     for (const change of changeLog) {
//       try {
//         await this.syncSingleChangeWithServer(change);
//         if (typeof change.changeId !== 'undefined') {
//           removeFromChangeLog(change.changeId);
//         }
//       } catch (error) {
//         console.error(`Failed to sync change ${change.changeId} with server:`, error);
//       }
//     }
//   },
//
//   async syncSingleChangeWithServer(change: ChangeLog): Promise<void> {
//     if (!change) {
//       throw new Error("Change is undefined");
//     }
//     switch (change.operation) {
//       case 'INSERT':
//       case 'UPDATE':
//         await this.syncWithServer(change.key, change.newValue);
//         break;
//       case 'DELETE':
//         await this.removeFromServer(change.key);
//         break;
//       default:
//         console.warn(`Unhandled change operation: ${change.operation}`);
//     }
//   },
//
//   async syncWithServer(key: keyof LocalStorage, value: LocalStorage[keyof LocalStorage]): Promise<void> {
//     try {
//       const url = `${API_BASE_URL}/api/sync/?key=${encodeURIComponent(key as string)}`;
//
//       const options: RequestInit = {
//         method: 'POST',
//         headers: {
//           'Content-Type': 'application/json',
//         },
//         mode: 'cors',
//         body: JSON.stringify({ value }), // We send only the value in the body as the key is already in the URL.
//       };
//
//       const response = await fetch(url, options);
//       const responseClone = response.clone(); // Clone the response to read it multiple times
//
//       if (!response.ok) {
//         let errorMessage = 'Server responded with an error';
//         try {
//           const errorBody = await responseClone.json(); // Try to parse as JSON first
//           errorMessage = errorBody.message || `Server responded with status code ${response.status}`;
//         } catch (jsonError) {
//           errorMessage = await responseClone.text(); // If response is not JSON, read as text
//         }
//
//         throw new Error(errorMessage);
//       }
//
//       // handle response logic, if needed
//       console.log('Data synced successfully with server');
//     } catch (error) {
//       console.error('Failed to sync with server:', error);
//     }
//   },
//
//
//
//
// const API_BASE_URL = "http://localhost:8080"
//
async function makeRequest(endpoint: string, method: string, body: any = null): Promise<any> {
  try {
    const url = `${API_BASE_URL}${endpoint}`;

    const options: RequestInit = {
      method,
      headers: {
        'Content-Type': 'application/json',
      },
      mode: 'cors',
    };
    if (body) {
      options.body = JSON.stringify(body);
    }

    const response = await fetch(url, options);

    if (!response.ok) {
      const errBody = await response.json();
      throw new Error(errBody.message || `Server responded with ${response.status}`);
    }

    return method === 'GET' ? response.json() : null;
  } catch (error) {
    console.error(`Request failed: ${error}`);
    return null;
  }
}
