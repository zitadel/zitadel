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
    GENERAL,
    IDP,
    LOCKOUT,
    LOGIN,
    LOGINTEXTS,
    MESSAGETEXTS,
    NOTIFICATIONS,
    OIDC,
    PRIVACYPOLICY,
    SECRETS,
} from '../../modules/settings-list/settings';

@Component({
  selector: 'cnsl-instance-settings',
  templateUrl: './instance-settings.component.html',
  styleUrls: ['./instance-settings.component.scss'],
})
export class InstanceSettingsComponent {
  public id: string = '';
  public PolicyComponentServiceType: any = PolicyComponentServiceType;
  public settingsList: SidenavSetting[] = [
    GENERAL,
    // notifications
    // { showWarn: true, ...NOTIFICATIONS },
    NOTIFICATIONS,
    // login
    LOGIN,
    COMPLEXITY,
    LOCKOUT,
    IDP,
    DOMAIN,
    // appearance
    BRANDING,
    MESSAGETEXTS,
    LOGINTEXTS,
    // others
    PRIVACYPOLICY,
    OIDC,
    SECRETS,
  ];
  constructor(breadcrumbService: BreadcrumbService, activatedRoute: ActivatedRoute) {
    const breadcrumbs = [
      new Breadcrumb({
        type: BreadcrumbType.INSTANCE,
        name: 'Instance',
        routerLink: ['/instance'],
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
