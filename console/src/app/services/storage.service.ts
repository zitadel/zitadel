import { Injectable } from '@angular/core';
import { OAuthStorage } from 'angular-oauth2-oidc';

const STORAGE_PREFIX = 'zitadel';

@Injectable({
    providedIn: 'root',
})
export class StorageService implements OAuthStorage {
    private storage: Storage = window.sessionStorage;

    constructor() { }

    public setItem<TValue = string>(key: string, value: TValue): void {
        this.storage.setItem(this.getPrefixedKey(key), JSON.stringify(value));
    }

    public getItem<TResult = string>(key: string): TResult | null {
        const result = this.storage.getItem(this.getPrefixedKey(key));
        if (result) {
            return JSON.parse(result);
        }
        return null;
    }

    public removeItem(key: string): void {
        this.storage.removeItem(this.getPrefixedKey(key));
    }

    public getPrefixedKey(key: string): string {
        return `${STORAGE_PREFIX}:${key}`;
    }
}

export class StorageConfig {
    clientId: string = '';
    storage: Storage = window.sessionStorage;
}

export enum StorageKey {
    organization = 'organization',
}
