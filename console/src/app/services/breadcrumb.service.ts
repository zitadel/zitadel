import { Injectable } from '@angular/core';
import { BehaviorSubject, switchMap } from 'rxjs';

import { ManagementService } from './mgmt.service';

export enum BreadcrumbType {
  INSTANCE,
  ORG,
  PROJECT,
  GRANTEDPROJECT,
  PROJECTGRANT,
  APP,
  IDP,
  AUTHUSER,
}

export class Breadcrumb {
  type: BreadcrumbType = BreadcrumbType.PROJECT;
  name?: string = '';
  param?: {
    key: 'projectid' | 'appid' | 'grantid' | 'id';
    value: string;
  } = {
    key: 'projectid',
    value: '',
  };
  routerLink: any[] = [];
  isZitadel?: boolean = false;
  hideNav?: boolean = false;

  constructor(init: Partial<Breadcrumb>) {
    Object.assign(this, init);
  }
}

@Injectable({
  providedIn: 'root',
})
export class BreadcrumbService {
  public readonly breadcrumbs$: BehaviorSubject<Breadcrumb[]> = new BehaviorSubject<Breadcrumb[]>([]);
  public readonly breadcrumbsExtended$ = this.breadcrumbs$.pipe(
    switchMap((breadcrumbs) => {
      const newValues = breadcrumbs.map(async (b) => {
        // the names of granted projects are set by the pages, since they can't be loaded by the project id only
        if (!b.name && b.type === BreadcrumbType.PROJECT && b.param) {
          b.name = await this.projectName(b.param.value);
        }
        return b;
      });
      return Promise.all(newValues);
    }),
  );

  constructor(private mgmtService: ManagementService) {}

  public setBreadcrumb(breadcrumbs: Breadcrumb[]) {
    this.breadcrumbs$.next(breadcrumbs);
  }

  private projectName(projectId: string): Promise<string> {
    return this.mgmtService
      .getProjectByID(projectId)
      .then((resp) => resp.project?.name ?? '')
      .catch(() => '');
  }
}
