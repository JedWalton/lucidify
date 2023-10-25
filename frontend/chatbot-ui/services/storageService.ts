// storageService.ts
import { LocalStorage } from '@/types/storage';
import { ChangeLog } from '@/types/changelog';


const CHANGE_LOG_KEY: keyof LocalStorage = '__CHANGE_LOG__';

function getChangeLog(): ChangeLog[] {
  const log = localStorage.getItem(CHANGE_LOG_KEY);
  return log ? JSON.parse(log) : [];
}

function addToChangeLog(change: Omit<ChangeLog, 'changeId'>): void {
  const log = getChangeLog();
  const lastChange = log[log.length - 1];
  const changeId = lastChange?.changeId ? lastChange.changeId + 1 : 1;

  log.push({ ...change, changeId });
  localStorage.setItem(CHANGE_LOG_KEY, JSON.stringify(log));
}

function clearChangeLog(): void {
  localStorage.removeItem(CHANGE_LOG_KEY);
}

function removeFromChangeLog(changeId: number): void {
  const log = getChangeLog();
  const index = log.findIndex(change => change.changeId === changeId);
  if (index !== -1) {
    log.splice(index, 1);
    localStorage.setItem(CHANGE_LOG_KEY, JSON.stringify(log));
  }
}

export const storageService = {
  async getItem(key: keyof LocalStorage): Promise<string | null> {
    let value = localStorage.getItem(key);

    if (value === null) {
      // If the value doesn't exist in local storage, try fetching from the server
      value = await this.fetchFromServer(key);
      if (value !== null) {
        // If found on the server, save it to local storage
        localStorage.setItem(key, JSON.parse(value).data);
        return value;
      }
    }
    return value;
  },

  async setItem(key: keyof LocalStorage, value: LocalStorage[keyof LocalStorage]): Promise<void> {
    const oldValue = await this.getItem(key);

    if (typeof value === 'object' && value !== null) {
      localStorage.setItem(key, JSON.stringify(value));
    } else {
      localStorage.setItem(key, String(value));
    }

    // Add to change log
    addToChangeLog({
      key,
      operation: oldValue ? 'UPDATE' : 'INSERT',
      oldValue: oldValue || '',
      newValue: value,
      timestamp: Date.now(),
    });
  },

  async removeItem(key: keyof LocalStorage): Promise<void> {
    const oldValue = await this.getItem(key);
    localStorage.removeItem(key);

    // Add to change log
    addToChangeLog({
      key,
      operation: 'DELETE',
      oldValue: oldValue || '',
      timestamp: Date.now(),
    });
  },

  async syncAllChangesWithServer(): Promise<void> {
    const changeLog = getChangeLog() as ChangeLog[];
    if (!changeLog || !changeLog.length) {
      throw new Error("ChangeLog is empty or not valid");
    }
    for (const change of changeLog) {
      try {
        await this.syncSingleChangeWithServer(change);
        if (typeof change.changeId !== 'undefined') {
          removeFromChangeLog(change.changeId);
        }
      } catch (error) {
        console.error(`Failed to sync change ${change.changeId} with server:`, error);
      }
    }
  },

  async syncSingleChangeWithServer(change: ChangeLog): Promise<void> {
    if (!change) {
      throw new Error("Change is undefined");
    }
    switch (change.operation) {
      case 'INSERT':
      case 'UPDATE':
        await this.syncWithServer(change.key, change.newValue);
        break;
      case 'DELETE':
        await this.removeFromServer(change.key);
        break;
      default:
        console.warn(`Unhandled change operation: ${change.operation}`);
    }
  },

  async syncWithServer(key: keyof LocalStorage, value: LocalStorage[keyof LocalStorage]): Promise<void> {
    try {
      const url = `${API_BASE_URL}/api/sync/?key=${encodeURIComponent(key as string)}`;

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
      console.log('Data synced successfully with server');
    } catch (error) {
      console.error('Failed to sync with server:', error);
    }
  },


  // async fetchFromServer(key: keyof LocalStorage): Promise<string | null> {
  //   const result = await makeRequest(`/api/sync/${key}`, 'GET'); // replace with your actual API endpoint
  // For GET and DELETE, it makes more sense to pass the 'key' in the query parameters (as you're not supposed to have a body in GET and DELETE requests).
  async fetchFromServer(key: keyof LocalStorage): Promise<string | null> {
    try {
      const url = `${API_BASE_URL}/api/sync/?key=${encodeURIComponent(key as string)}`;

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

  async removeFromServer(key: keyof LocalStorage): Promise<void> {
    const result = await makeRequest(`/api/sync?key=${key}`, 'DELETE'); // Sending 'key' in query parameters
    if (result !== null) {
      console.log(`Data associated with '${key}' successfully deleted from the server.`);
    }
  }
};

const API_BASE_URL = "http://localhost:8080"

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
