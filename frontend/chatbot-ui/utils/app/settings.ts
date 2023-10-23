import { storageService } from '@/services/storageService';
import { Settings } from '@/types/settings';

const STORAGE_KEY = 'settings';

export const getSettings = (): Settings => {
  let settings: Settings = {
    theme: 'dark',
  };
  const settingsJson = storageService.getItem(STORAGE_KEY);
  if (settingsJson) {
    try {
      let savedSettings = JSON.parse(settingsJson) as Settings;
      settings = Object.assign(settings, savedSettings);
    } catch (e) {
      console.error(e);
    }
  }
  return settings;
};

export const saveSettings = (settings: Settings) => {
  storageService.setItem(STORAGE_KEY, JSON.stringify(settings));
};
