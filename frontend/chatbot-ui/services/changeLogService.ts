// changeLogService.ts
import { LocalStorage } from '@/types/storage';
import { ChangeLog } from '@/types/changelog';

const CHANGE_LOG_KEY: keyof LocalStorage = '__CHANGE_LOG__';

export const changeLogService = {
  getChangeLog(): ChangeLog[] {
    const log = localStorage.getItem(CHANGE_LOG_KEY);
    return log ? JSON.parse(log) : [];
  },

  addToChangeLog(change: Omit<ChangeLog, 'changeId'>): void {
    const log = this.getChangeLog();
    const lastChange = log[log.length - 1];
    const changeId = lastChange?.changeId ? lastChange.changeId + 1 : 1;

    log.push({ ...change, changeId });
    localStorage.setItem(CHANGE_LOG_KEY, JSON.stringify(log));
  },

  clearChangeLog(): void {
    localStorage.removeItem(CHANGE_LOG_KEY);
  },

  removeFromChangeLog(changeId: number): void {
    const log = this.getChangeLog();
    const index = log.findIndex(change => change.changeId === changeId);
    if (index !== -1) {
      log.splice(index, 1);
      localStorage.setItem(CHANGE_LOG_KEY, JSON.stringify(log));
    }
  },

  async syncChangeLogToServer() {
    const changelog = changeLogService.getChangeLog();

    const response = await fetch(process.env.PUBLIC_BACKEND_API_URL + '/api/sync/changelog', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify(changelog)
    });

    if (!response.ok) {
        // Handle error
        const data = await response.json();
        console.error(data.message);
    }
}
}

