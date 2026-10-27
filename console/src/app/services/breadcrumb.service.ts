import { Injectable } from '@angular/core';
import { BehaviorSubject, concat, from, of, switchMap } from 'rxjs';

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
      // the names of granted projects are set by the pages, since they can't be loaded by the project id only
      const unnamed = breadcrumbs.filter((b) => !b.name && b.type === BreadcrumbType.PROJECT && b.param);
      if (!unnamed.length) {
        return of(breadcrumbs);
      }
      // show the breadcrumbs right away and update them once the names are loaded
      const named = Promise.all(
        unnamed.map(async (b) => {
          b.name = await this.projectName(b.param?.value ?? '');
        }),
      ).then(() => [...breadcrumbs]);
      return concat(of(breadcrumbs), from(named));
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
