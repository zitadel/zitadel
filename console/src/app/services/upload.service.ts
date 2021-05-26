import { HttpClient } from '@angular/common/http';
import { Injectable } from '@angular/core';

import { Org } from '../proto/generated/zitadel/org_pb';
import { StorageService } from './storage.service';

const ORG_STORAGE_KEY = 'organization';
const authorizationKey = 'Authorization';
const orgKey = 'x-zitadel-orgid';

const bearerPrefix = 'Bearer';
const accessTokenStorageKey = 'access_token';

export enum UploadEndpoint {
  IAMDARKLOGO = 'iam/policy/label/logo/dark',
  IAMLIGHTLOGO = 'iam/policy/label/logo/light',
  IAMDARKICON = 'iam/policy/label/icon/dark',
  IAMLIGHTICON = 'iam/policy/label/icon/light',

  MGMTDARKLOGO = 'management/policy/label/logo/dark',
  MGMTLIGHTLOGO = 'management/policy/label/logo/light',
  MGMTDARKICON = 'management/policy/label/icon/dark',
  MGMTLIGHTICON = 'management/policy/label/icon/light',
}

@Injectable({
  providedIn: 'root',
})
export class UploadService {
  private serviceUrl: string = '';
  private accessToken: string = '';
  private org!: Org.AsObject;
  constructor(private http: HttpClient, private storageService: StorageService) {

    http.get('./assets/environment.json')
      .toPromise().then((data: any) => {
        if (data && data.uploadServiceUrl) {
          this.serviceUrl = data.uploadServiceUrl;
          const aT = this.storageService.getItem(accessTokenStorageKey);

          if (aT) {
            this.accessToken = aT;
          }

          const org: Org.AsObject | null = (this.storageService.getItem(ORG_STORAGE_KEY));

          if (org) {
            this.org = org;
          }
        }
      }).catch(error => {
        console.error(error);
      });
  }

  public upload(endpoint: UploadEndpoint, body: any): Promise<any> {
    return this.http.post(`${this.serviceUrl}/assets/v1/${endpoint}`,
      body,
      {
        headers: {
          [authorizationKey]: `${bearerPrefix} ${this.accessToken}`,
          [orgKey]: `${this.org.id}`,
        },
      }).toPromise();
  }
}
