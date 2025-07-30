import { Injectable } from '@angular/core';
import { BehaviorSubject, combineLatest, map } from 'rxjs';

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
  public readonly breadcrumbsExtended$ = combineLatest([
    this.breadcrumbs$,
    this.mgmtService.ownedProjects,
    this.mgmtService.grantedProjects,
  ]).pipe(
    map(([breadcrumbs, projects, grantedProjects]) => {
      const newValues = breadcrumbs.map((b) => {
        if (!b.name && b.type === BreadcrumbType.PROJECT) {
          const project = projects.find((project) => b.param && project.id === b.param.value);
          b.name = project?.name ?? '';
          return b;
        } else if (!b.name && b.type === BreadcrumbType.GRANTEDPROJECT) {
          const grantedproject = grantedProjects.find(
            (grantedproject) => b.param && grantedproject.projectId === b.param.value,
          );
          b.name = grantedproject?.projectName ?? '';
          return b;
        } else {
          return b;
        }
      });
      return newValues;
    }),
  );

  constructor(private mgmtService: ManagementService) {}

  public setBreadcrumb(breadcrumbs: Breadcrumb[]) {
    this.breadcrumbs$.next(breadcrumbs);
  }
}
