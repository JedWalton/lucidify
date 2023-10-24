// storageService.ts
import { LocalStorage } from '@/types/storage';
import { Conversation } from '@/types/chat';
import { FolderInterface } from '@/types/folder';
import { PluginKey } from '@/types/plugin';
import { Prompt } from '@/types/prompt';
import { Settings } from '@/types/settings';

export const storageService = {
  async getItem(key: keyof LocalStorage): Promise<string | null> {
    let value = localStorage.getItem(key);

    if (value === null) {
      // If the value doesn't exist in local storage, try fetching from the server
      value = await this.fetchFromServer(key);
      if (value !== null) {
        // If found on the server, save it to local storage
        localStorage.setItem(key, value);
      }
    }
    return value;
  },

  async setItem(key: keyof LocalStorage, value: LocalStorage[keyof LocalStorage]): Promise<void> {
    // For complex types, we need to stringify them to store in localStorage
    if (typeof value === 'object' && value !== null) {
      localStorage.setItem(key, JSON.stringify(value));
    } else {
      localStorage.setItem(key, String(value));
    }
        // Also sync with server
    await this.syncWithServer(key, value);
  },

  async removeItem(key: keyof LocalStorage): Promise<void> {
    localStorage.removeItem(key);
    // Also attempt to remove the data from the server
    await this.removeFromServer(key).catch(error => {
      console.error('Failed to remove item from server:', error);
    });
  },

  // Now, refactor your methods using the above helper function
  async syncWithServer(key: keyof LocalStorage, value: LocalStorage[keyof LocalStorage]): Promise<void> {
    const body = { key, value };
    const result = await makeRequest('/api/sync', 'POST', body); // replace '/api/sync' with your actual API endpoint
    if (result === null) {
      console.error('Failed to sync with server');
    } else {
      // handle response logic, if needed
    }
  },

  async fetchFromServer(key: keyof LocalStorage): Promise<string | null> {
    const result = await makeRequest(`/api/sync/${key}`, 'GET'); // replace with your actual API endpoint
    if (result === null) {
      console.error('Failed to fetch data from server');
      return null;
    }

    // If the server returns the data directly, we stringify it to keep the method's signature consistent
    // You might need to adjust this based on how your server responds
    return JSON.stringify(result);
  },

  async removeFromServer(key: keyof LocalStorage): Promise<void> {
    const result = await makeRequest(`/api/sync/${key}`, 'DELETE');
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
    };
    if (body) {
      options.body = JSON.stringify(body);
    }

    const response = await fetch(url, options);

    if (!response.ok) {
      throw new Error(`Server responded with ${response.status}`);
    }

    return method === 'GET' ? response.json() : null;
  } catch (error) {
    console.error(`Request failed: ${error}`);
    return null;
  }
}
