import { HttpClient } from '@angular/common/http';
import { Injectable } from '@angular/core';

import { StorageService } from './storage.service';

const authorizationKey = 'Authorization';
const bearerPrefix = 'Bearer';
const accessTokenStorageKey = 'access_token';

export interface StripeCustomer {
  contact: string;
  company?: string;
  address: string;
  city: string;
  postal_code: string;
  country: string;
}

@Injectable({
  providedIn: 'root',
})
export class UploadService {
  constructor(private http: HttpClient, private storageService: StorageService) { }

  public getLink(orgId: string, redirectURI: string): Promise<any> {
    return this.http.get('./assets/environment.json')
      .toPromise().then((data: any) => {
        if (data && data.subscriptionServiceUrl) {
          const serviceUrl = data.subscriptionServiceUrl;
          const accessToken = this.storageService.getItem(accessTokenStorageKey);
          return this.http.get(`${serviceUrl}/redirect`, {
            headers: {
              [authorizationKey]: `${bearerPrefix} ${accessToken}`,
            },
            params: {
              'org': orgId,
              'return_url': encodeURI(redirectURI),
              'country': 'ch',
            },
          }).toPromise();
        } else {
          return Promise.reject('Could not load environment');
        }
      });
  }

  public getCustomer(orgId: string): Promise<any> {
    return this.http.get('./assets/environment.json')
      .toPromise().then((data: any) => {
        // if (data && data.subscriptionServiceUrl) {
        //   const serviceUrl = data.subscriptionServiceUrl;
        //   const accessToken = this.storageService.getItem(accessTokenStorageKey);
        //   return this.http.get(`${serviceUrl}/customer`, {
        //     headers: {
        //       [authorizationKey]: `${bearerPrefix} ${accessToken}`,
        //     },
        //     params: {
        //       'org': orgId,
        //     },
        //   }).toPromise();
        // } else {
        //   return Promise.reject('Could not load environment');
        // }
      });
  }

  public setCustomer(orgId: string, body: StripeCustomer): Promise<any> {
    return this.http.get('./assets/environment.json')
      .toPromise().then((data: any) => {
        // if (data && data.subscriptionServiceUrl) {
        //   const serviceUrl = data.subscriptionServiceUrl;
        //   const accessToken = this.storageService.getItem(accessTokenStorageKey);
        //   return this.http.post(`${serviceUrl}/customer`, body, {
        //     headers: {
        //       [authorizationKey]: `${bearerPrefix} ${accessToken}`,
        //     },
        //     params: {
        //       'org': orgId,
        //     },
        //   }).toPromise();
        // } else {
        //   return Promise.reject('Could not load environment');
        // }
      });
  }
}
