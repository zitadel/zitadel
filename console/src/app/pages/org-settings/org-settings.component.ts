import { Component, OnInit } from '@angular/core';
import { ActivatedRoute, Params } from '@angular/router';
import { forkJoin, of, take } from 'rxjs';
import { PolicyComponentServiceType } from 'src/app/modules/policies/policy-component-types.enum';
import { SidenavSetting } from 'src/app/modules/sidenav/sidenav.component';
import { Breadcrumb, BreadcrumbService, BreadcrumbType } from 'src/app/services/breadcrumb.service';

import { GrpcAuthService } from 'src/app/services/grpc-auth.service';
import {
  BRANDING,
  COMPLEXITY,
  DOMAIN,
  IDP,
  LOCKOUT,
  LOGIN,
  LOGINTEXTS,
  MESSAGETEXTS,
  NOTIFICATION_POLICY,
  PRIVACYPOLICY,
} from '../../modules/settings-list/settings';

@Component({
  selector: 'cnsl-org-settings',
  templateUrl: './org-settings.component.html',
  styleUrls: ['./org-settings.component.scss'],
})
export class OrgSettingsComponent implements OnInit {
  public id: string = '';
  public PolicyComponentServiceType: any = PolicyComponentServiceType;

  private defaultSettingsList: SidenavSetting[] = [
    LOGIN,
    IDP,
    COMPLEXITY,
    LOCKOUT,
    NOTIFICATION_POLICY,
    DOMAIN,
    BRANDING,
    MESSAGETEXTS,
    LOGINTEXTS,
    PRIVACYPOLICY,
  ];

  public settingsList: SidenavSetting[] = [];

  constructor(
    breadcrumbService: BreadcrumbService,
    activatedRoute: ActivatedRoute,
    public authService: GrpcAuthService,
  ) {
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

  ngOnInit(): void {
    checkSettingsPermissions(this.defaultSettingsList, PolicyComponentServiceType.MGMT, this.authService).subscribe(
      (allowed) => {
        this.settingsList = this.defaultSettingsList.filter((setting, index) => {
          return allowed[index];
        });
      },
    );
  }
}

// Return a Observables<boolean>[] that will wait till all service calls are finished to then check if user is allowed to see a setting
export function checkSettingsPermissions(settings: SidenavSetting[], serviceType: string, authService: GrpcAuthService) {
  return forkJoin(
    settings
      .filter((setting) => {
        if (serviceType === PolicyComponentServiceType.ADMIN) {
          return setting.requiredRoles && setting.requiredRoles.admin;
        } else {
          return setting.requiredRoles && setting.requiredRoles.mgmt;
        }
      })
      .map((setting) => {
        if (!setting.requiredRoles) {
          return of(true);
        }

        if (!setting.requiredRoles.mgmt) {
          return of(true);
        }

        if (setting.requiredRoles.mgmt) {
          return authService.isAllowed(setting.requiredRoles.mgmt).pipe(take(1));
        }
        return of(false);
      }),
  );
}
