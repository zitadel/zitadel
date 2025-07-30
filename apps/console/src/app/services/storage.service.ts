import { Injectable } from '@angular/core';
import { OAuthStorage } from 'angular-oauth2-oidc';

const STORAGE_PREFIX = 'zitadel';

@Injectable({
  providedIn: 'root',
})
export class StorageService implements OAuthStorage {
  private sessionStorage: Storage = window.sessionStorage;
  private localStorage: Storage = window.localStorage;

  constructor() {}

  public setItem<TValue = string>(key: string, value: TValue, location: StorageLocation = StorageLocation.session): void {
    this.getStorage(location).setItem(this.getPrefixedKey(key), JSON.stringify(value));
  }

  public getItem<TResult = string>(key: string, location: StorageLocation = StorageLocation.session): TResult | null {
    const result = this.getStorage(location).getItem(this.getPrefixedKey(key));
    if (result) {
      return JSON.parse(result);
    }
    return null;
  }

  public removeItem(key: string, location: StorageLocation = StorageLocation.session): void {
    this.getStorage(location).removeItem(this.getPrefixedKey(key));
  }

  public getPrefixedKey(key: string): string {
    return `${STORAGE_PREFIX}:${key}`;
  }

  private getStorage(location: StorageLocation): Storage {
    return location === StorageLocation.session ? this.sessionStorage : this.localStorage;
  }
}

export class StorageConfig {
  clientId: string = '';
  storage: Storage = window.sessionStorage;
}

export enum StorageKey {
  organization = 'organization',
}

export enum StorageLocation {
  session,
  local,
}
