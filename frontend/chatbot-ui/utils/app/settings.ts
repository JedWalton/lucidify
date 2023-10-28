import { Settings } from '@/types/settings';

const STORAGE_KEY = 'settings';

export const getSettings = async (): Promise<Settings> => {
  let settings: Settings = {
    theme: 'dark',
  };

  try {
    const settingsJson = localStorage.getItem(STORAGE_KEY); // await the Promise
    if (settingsJson) {
      const savedSettings = JSON.parse(settingsJson) as Settings;
      settings = Object.assign(settings, savedSettings);
    }
  } catch (e) {
    console.error(e);
  }

  return settings;
};

export const saveSettings = async (settings: Settings) => {
  try {
    localStorage.setItem(STORAGE_KEY, JSON.stringify(settings));
  } catch (e) {
    console.error('Failed to save settings:', e);
  }
};

