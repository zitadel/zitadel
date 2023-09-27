import { Component, OnDestroy, OnInit } from '@angular/core';
import { ActivatedRoute, Params } from '@angular/router';
import { Subject, takeUntil } from 'rxjs';
import { PolicyComponentServiceType } from 'src/app/modules/policies/policy-component-types.enum';
import { SidenavSetting } from 'src/app/modules/sidenav/sidenav.component';
import { Breadcrumb, BreadcrumbService, BreadcrumbType } from 'src/app/services/breadcrumb.service';

import { GrpcAuthService } from 'src/app/services/grpc-auth.service';
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
  SECURITY,
} from '../../modules/settings-list/settings';
import { checkSettingsPermissions } from '../org-settings/org-settings.component';

@Component({
  selector: 'cnsl-instance-settings',
  templateUrl: './instance-settings.component.html',
  styleUrls: ['./instance-settings.component.scss'],
})
export class InstanceSettingsComponent implements OnInit, OnDestroy {
  public id: string = '';
  public PolicyComponentServiceType: any = PolicyComponentServiceType;
  public defaultSettingsList: SidenavSetting[] = [
    GENERAL,
    // notifications
    // { showWarn: true, ...NOTIFICATIONS },
    NOTIFICATIONS,
    // login
    LOGIN,
    IDP,
    COMPLEXITY,
    LOCKOUT,

    DOMAIN,
    // appearance
    BRANDING,
    MESSAGETEXTS,
    LOGINTEXTS,
    // others
    PRIVACYPOLICY,
    OIDC,
    SECRETS,
    SECURITY,
  ];

  public settingsList: SidenavSetting[] = [];

  private destroy$: Subject<void> = new Subject();
  constructor(breadcrumbService: BreadcrumbService, activatedRoute: ActivatedRoute, public authService: GrpcAuthService) {
    const breadcrumbs = [
      new Breadcrumb({
        type: BreadcrumbType.INSTANCE,
        name: 'Instance',
        routerLink: ['/instance'],
      }),
    ];
    breadcrumbService.setBreadcrumb(breadcrumbs);

    activatedRoute.queryParams.pipe(takeUntil(this.destroy$)).subscribe((params: Params) => {
      const { id } = params;
      if (id) {
        this.id = id;
      }
    });
  }

  ngOnInit(): void {
    checkSettingsPermissions(this.defaultSettingsList, PolicyComponentServiceType.ADMIN, this.authService).subscribe(
      (allowed) => {
        this.settingsList = this.defaultSettingsList.filter((setting, index) => {
          return allowed[index];
        });
      },
    );
  }

  ngOnDestroy(): void {
    this.destroy$.next();
    this.destroy$.complete();
  }
}
