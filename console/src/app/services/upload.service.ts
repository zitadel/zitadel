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
  IAMFONT = 'iam/policy/label/font',
  MGMTFONT = 'org/policy/label/font',

  IAMDARKLOGO = 'iam/policy/label/logo/dark',
  IAMLIGHTLOGO = 'iam/policy/label/logo',
  IAMDARKICON = 'iam/policy/label/icon/dark',
  IAMLIGHTICON = 'iam/policy/label/icon',

  MGMTDARKLOGO = 'org/policy/label/logo/dark',
  MGMTLIGHTLOGO = 'org/policy/label/logo',
  MGMTDARKICON = 'org/policy/label/icon/dark',
  MGMTLIGHTICON = 'org/policy/label/icon',
}

export enum DownloadEndpoint {
  IAMFONT = 'iam/policy/label/font',
  MGMTFONT = 'org/policy/label/font',

  IAMDARKLOGO = 'iam/policy/label/logo/dark',
  IAMLOGO = 'iam/policy/label/logo',
  IAMDARKICON = 'iam/policy/label/icon/dark',
  IAMICON = 'iam/policy/label/icon',

  MGMTDARKLOGO = 'org/policy/label/logo/dark',
  MGMTLOGO = 'org/policy/label/logo',
  MGMTDARKICON = 'org/policy/label/icon/dark',
  MGMTICON = 'org/policy/label/icon',

  IAMDARKLOGOPREVIEW = 'iam/policy/label/logo/dark/_preview',
  IAMLOGOPREVIEW = 'iam/policy/label/logo/_preview',
  IAMDARKICONPREVIEW = 'iam/policy/label/icon/dark/_preview',
  IAMICONPREVIEW = 'iam/policy/label/icon/_preview',

  MGMTDARKLOGOPREVIEW = 'org/policy/label/logo/dark/_preview',
  MGMTLOGOPREVIEW = 'org/policy/label/logo/_preview',
  MGMTDARKICONPREVIEW = 'org/policy/label/icon/dark/_preview',
  MGMTICONPREVIEW = 'org/policy/label/icon/_preview',
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

  public upload(endpoint: UploadEndpoint | string, body: any): Promise<any> {
    return this.http.post(`${this.serviceUrl}/assets/v1/${endpoint}`,
      body,
      {
        headers: {
          [authorizationKey]: `${bearerPrefix} ${this.accessToken}`,
          [orgKey]: `${this.org.id}`,
        },
      }).toPromise();
  }

  public load(endpoint: string): Promise<any> {
    return this.http.get(`${this.serviceUrl}/assets/v1/${endpoint}`,

      {
        responseType: 'blob',
        headers: {
          [authorizationKey]: `${bearerPrefix} ${this.accessToken}`,
          [orgKey]: `${this.org.id}`,
        },
      }).toPromise();
  }
}
