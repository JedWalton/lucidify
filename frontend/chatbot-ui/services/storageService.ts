// storageService.ts
import { LocalStorage } from '@/types/storage';
import { Conversation } from '@/types/chat';
import { FolderInterface } from '@/types/folder';
import { PluginKey } from '@/types/plugin';
import { Prompt } from '@/types/prompt';
import { Settings } from '@/types/settings';

// import { LocalStorage, Conversation, FolderInterface, PluginKey, Prompt } from './your-types-file'; // update with your actual file path

export const storageService = {
  getItem(key: keyof LocalStorage): string | null {
    const value = localStorage.getItem(key);
    return value; // return the string directly, parse after retrieval where necessary
  },
  setItem(key: keyof LocalStorage, value: LocalStorage[keyof LocalStorage]): void {
    // For complex types, we need to stringify them to store in localStorage
    if (typeof value === 'object' && value !== null) {
      localStorage.setItem(key, JSON.stringify(value));
    } else {
      localStorage.setItem(key, String(value));
    }
  },
  removeItem(key: keyof LocalStorage): void {
    localStorage.removeItem(key);
  },
};

