import { Injectable } from '@angular/core';
import { BehaviorSubject } from 'rxjs';

export enum BreadcrumbType {
  PROJECT,
  GRANTEDPROJECT,
  APP,
  IDP,
}

export class Breadcrumb {
  type: BreadcrumbType = BreadcrumbType.PROJECT;
  name: string = '';
  param: {
    key: 'projectid' | 'appid' | 'id';
    value: string;
  } = {
    key: 'projectid',
    value: '',
  };
  routerLink: any[] = [];

  constructor(init: Partial<Breadcrumb>) {
    Object.assign(this, init);
  }
}

@Injectable({
  providedIn: 'root',
})
export class BreadcrumbService {
  public readonly breadcrumbs$: BehaviorSubject<Breadcrumb[]> = new BehaviorSubject<Breadcrumb[]>([]);
  constructor() {}

  public setBreadcrumb(breadcrumbs: Breadcrumb[]) {
    this.breadcrumbs$.next(breadcrumbs);
  }
}
