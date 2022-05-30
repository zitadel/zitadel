import { Component } from '@angular/core';
import { ActivatedRoute, Params } from '@angular/router';
import { take } from 'rxjs';
import { PolicyComponentServiceType } from 'src/app/modules/policies/policy-component-types.enum';
import { SidenavSetting } from 'src/app/modules/sidenav/sidenav.component';
import { Breadcrumb, BreadcrumbService, BreadcrumbType } from 'src/app/services/breadcrumb.service';

import {
    BRANDING,
    COMPLEXITY,
    DOMAIN,
    IDP,
    LOCKOUT,
    LOGIN,
    LOGINTEXTS,
    MESSAGETEXTS,
    PRIVACYPOLICY,
} from '../../modules/settings-list/settings';

@Component({
  selector: 'cnsl-org-settings',
  templateUrl: './org-settings.component.html',
  styleUrls: ['./org-settings.component.scss'],
})
export class OrgSettingsComponent {
  public id: string = '';
  public PolicyComponentServiceType: any = PolicyComponentServiceType;
  public settingsList: SidenavSetting[] = [
    LOGIN,
    COMPLEXITY,
    LOCKOUT,
    IDP,
    DOMAIN,
    BRANDING,
    MESSAGETEXTS,
    LOGINTEXTS,
    PRIVACYPOLICY,
  ];

  constructor(breadcrumbService: BreadcrumbService, activatedRoute: ActivatedRoute) {
    const breadcrumbs = [
      new Breadcrumb({
        type: BreadcrumbType.ORG,
        routerLink: ['/org'],
      }),
    ];
    breadcrumbService.setBreadcrumb(breadcrumbs);

    activatedRoute.queryParams.pipe(take(1)).subscribe((params: Params) => {
      const { id } = params;
      if (id) {
        this.id = id;
      }
    });
  }
}
