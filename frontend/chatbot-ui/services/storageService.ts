// storageService.ts
import { LocalStorage } from '@/types/storage';
import { Conversation } from '../types/chat';
import { FolderInterface } from '../types/folder';
import { PluginKey } from '../types/plugin';
import { Prompt } from '../types/prompt';

// import { LocalStorage, Conversation, FolderInterface, PluginKey, Prompt } from './your-types-file'; // update with your actual file path

export const storageService = {
  get(key: keyof LocalStorage): LocalStorage[keyof LocalStorage] | null {
    const value = localStorage.getItem(key);
    if (!value) return null;

    // For complex objects stored as strings, we'll need to parse them
    try {
      return JSON.parse(value);
    } catch {
      // If an error occurs, it's possible that the item is a string that isn't JSON, so return it directly
      return value;
    }
  },
  set(key: keyof LocalStorage, value: LocalStorage[keyof LocalStorage]): void {
    // For complex types, we need to stringify them to store in localStorage
    if (typeof value === 'object' && value !== null) {
      localStorage.setItem(key, JSON.stringify(value));
    } else {
      localStorage.setItem(key, String(value));
    }
  },
  remove(key: keyof LocalStorage): void {
    localStorage.removeItem(key);
  },
};

