import { HttpClient } from '@angular/common/http';
import { Injectable } from '@angular/core';

import { StorageService } from './storage.service';



const authorizationKey = 'Authorization';
const bearerPrefix = 'Bearer';
const accessTokenStorageKey = 'access_token';

@Injectable({
    providedIn: 'root',
})
export class SubscriptionService {
    constructor(private http: HttpClient, private storageService: StorageService) { }

    public getLink(orgId: string, redirectURI: string): Promise<any> {
        return this.http.get('./assets/environment.json')
            .toPromise().then((data: any) => {
                if (data && data.subscriptionServiceUrl) {
                    const serviceUrl = data.subscriptionServiceUrl;
                    console.log(serviceUrl);

                    const accessToken = this.storageService.getItem(accessTokenStorageKey);
                    console.log(accessToken);

                    return this.http.get(serviceUrl, {
                        headers: {
                            // 'Content-Type': 'application/json; charset=utf-8',
                            [authorizationKey]: `${bearerPrefix} ${accessToken}`
                        },
                        params: {
                            'org': orgId,
                            'return_url': encodeURI(redirectURI),
                            'country': 'ch'
                        }
                    }).toPromise();
                } else {
                    return Promise.reject('Could not load environment');
                }
            });
    }
}