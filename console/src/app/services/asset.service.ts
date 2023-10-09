import { HttpClient } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { switchMap } from 'rxjs';

import { PolicyComponentServiceType } from '../modules/policies/policy-component-types.enum';
import { Theme } from '../modules/policies/private-labeling-policy/private-labeling-policy.component';
import { EnvironmentService } from './environment.service';
import { StorageService } from './storage.service';

const authorizationKey = 'Authorization';
const orgKey = 'x-zitadel-orgid';

const bearerPrefix = 'Bearer';
const accessTokenStorageKey = 'access_token';

export enum AssetType {
  LOGO,
  ICON,
}

export enum AssetEndpoint {
  IAMFONT = 'instance/policy/label/font',
  MGMTFONT = 'org/policy/label/font',

  IAMDARKLOGO = 'instance/policy/label/logo/dark',
  IAMLOGO = 'instance/policy/label/logo',
  IAMDARKICON = 'instance/policy/label/icon/dark',
  IAMICON = 'instance/policy/label/icon',

  MGMTDARKLOGO = 'org/policy/label/logo/dark',
  MGMTLOGO = 'org/policy/label/logo',
  MGMTDARKICON = 'org/policy/label/icon/dark',
  MGMTICON = 'org/policy/label/icon',

  IAMDARKLOGOPREVIEW = 'instance/policy/label/logo/dark/_preview',
  IAMLOGOPREVIEW = 'instance/policy/label/logo/_preview',
  IAMDARKICONPREVIEW = 'instance/policy/label/icon/dark/_preview',
  IAMICONPREVIEW = 'instance/policy/label/icon/_preview',

  MGMTDARKLOGOPREVIEW = 'org/policy/label/logo/dark/_preview',
  MGMTLOGOPREVIEW = 'org/policy/label/logo/_preview',
  MGMTDARKICONPREVIEW = 'org/policy/label/icon/dark/_preview',
  MGMTICONPREVIEW = 'org/policy/label/icon/_preview',
}

export const ENDPOINT = {
  [Theme.DARK]: {
    [PolicyComponentServiceType.ADMIN]: {
      [AssetType.LOGO]: AssetEndpoint.IAMDARKLOGO,
      [AssetType.ICON]: AssetEndpoint.IAMDARKICON,
    },
    [PolicyComponentServiceType.MGMT]: {
      [AssetType.LOGO]: AssetEndpoint.MGMTDARKLOGO,
      [AssetType.ICON]: AssetEndpoint.MGMTDARKICON,
    },
  },
  [Theme.LIGHT]: {
    [PolicyComponentServiceType.ADMIN]: {
      [AssetType.LOGO]: AssetEndpoint.IAMLOGO,
      [AssetType.ICON]: AssetEndpoint.IAMICON,
    },
    [PolicyComponentServiceType.MGMT]: {
      [AssetType.LOGO]: AssetEndpoint.MGMTLOGO,
      [AssetType.ICON]: AssetEndpoint.MGMTICON,
    },
  },
};

@Injectable({
  providedIn: 'root',
})
export class AssetService {
  private accessToken: string = '';
  constructor(private envService: EnvironmentService, private http: HttpClient, private storageService: StorageService) {
    const aT = this.storageService.getItem(accessTokenStorageKey);
    if (aT) {
      this.accessToken = aT;
    }
  }

  public upload(endpoint: AssetEndpoint | string, body: any, orgId?: string): Promise<any> {
    const headers: any = {
      [authorizationKey]: `${bearerPrefix} ${this.accessToken}`,
    };
    if (orgId) {
      headers[orgKey] = `${orgId}`;
    }
    return this.envService.env
      .pipe(
        switchMap((env) =>
          this.http.post(`${env.api}/assets/v1/${endpoint}`, body, {
            headers: headers,
          }),
        ),
      )
      .toPromise();
  }
}
