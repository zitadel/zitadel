import { Injectable } from '@angular/core';
import { OAuthStorage } from 'angular-oauth2-oidc';

import { GrpcService } from './grpc.service';


@Injectable({
    providedIn: 'root',
})

export class StorageService implements OAuthStorage {
    private storage: Storage = window.sessionStorage;

    constructor(private grpcService: GrpcService) { }

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
        return `caos:${key}`;
    }
}

export class StorageConfig {
    clientId: string = '';
    storage: Storage = window.sessionStorage;
}

export enum StorageKey {
    organization = 'organization',
}
