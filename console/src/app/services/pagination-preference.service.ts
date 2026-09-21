import { inject, Injectable } from '@angular/core';
import { StorageLocation, StorageService } from './storage.service';

const KEY_PREFIX = 'pagesize-';

/**
 * Remembers the page size a user picked per table, so it survives navigation and reloads.
 *
 * Stored values are validated against the options the table actually offers. That guards
 * against keys written by an older option set as well as hand edited storage, either of
 * which would otherwise end up in a request the backend rejects.
 */
@Injectable({
  providedIn: 'root',
})
export class PaginationPreferenceService {
  private readonly storage = inject(StorageService);

  public get(key: string, fallback: number, allowedOptions?: Array<number>): number {
    let stored: number | null = null;
    try {
      stored = this.storage.getItem<number>(KEY_PREFIX + key, StorageLocation.local);
    } catch {
      // storage can be unavailable (private mode, blocked by the browser) or hold invalid JSON
      return fallback;
    }

    if (typeof stored !== 'number' || !Number.isInteger(stored) || stored <= 0) {
      return fallback;
    }
    if (allowedOptions?.length && !allowedOptions.includes(stored)) {
      return fallback;
    }
    return stored;
  }

  public set(key: string, size: number): void {
    try {
      this.storage.setItem<number>(KEY_PREFIX + key, size, StorageLocation.local);
    } catch {
      // remembering the page size is a convenience and must never break the table
    }
  }
}
