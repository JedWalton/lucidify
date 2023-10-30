import { storageService } from '@/services/storageService';
import { FolderInterface } from '@/types/folder';

export const saveFolders = async (folders: FolderInterface[]) => {
  localStorage.setItem('folders', JSON.stringify(folders));
  await storageService.setItemOnServer('folders', JSON.stringify(folders));
};
