import { Injectable } from '@angular/core';
import { OAuthStorage } from 'angular-oauth2-oidc';

@Injectable({ providedIn: 'root' })
export class AccountStorage implements OAuthStorage {
  getItem(key: string): string | null {
    return sessionStorage.getItem(`account:v1:${key}`);
  }
  setItem(key: string, value: string): void {
    sessionStorage.setItem(`account:v1:${key}`, value);
  }
  removeItem(key: string): void {
    sessionStorage.removeItem(`account:v1:${key}`);
  }
  clear(): void {
    for (const key of Object.keys(sessionStorage)) if (key.startsWith('account:v1:')) sessionStorage.removeItem(key);
  }
}
