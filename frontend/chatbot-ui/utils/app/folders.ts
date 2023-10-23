import { storageService } from '@/services/storageService';
import { FolderInterface } from '@/types/folder';

export const saveFolders = (folders: FolderInterface[]) => {
  storageService.setItem('folders', JSON.stringify(folders));
};
