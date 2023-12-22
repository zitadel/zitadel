import { Component, OnDestroy, OnInit } from '@angular/core';
import { ActivatedRoute, Params } from '@angular/router';
import { Observable, of, Subject, takeUntil } from 'rxjs';
import { PolicyComponentServiceType } from 'src/app/modules/policies/policy-component-types.enum';
import { SidenavSetting } from 'src/app/modules/sidenav/sidenav.component';
import { Breadcrumb, BreadcrumbService, BreadcrumbType } from 'src/app/services/breadcrumb.service';

import { GrpcAuthService } from 'src/app/services/grpc-auth.service';
import {
  BRANDING,
  COMPLEXITY,
  DOMAIN,
  LANGUAGES,
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
  SMS_PROVIDER,
  SMTP_PROVIDER,
} from '../../modules/settings-list/settings';

@Component({
  selector: 'cnsl-instance-settings',
  templateUrl: './instance-settings.component.html',
  styleUrls: ['./instance-settings.component.scss'],
})
export class InstanceSettingsComponent implements OnInit, OnDestroy {
  public id: string = '';
  public PolicyComponentServiceType: any = PolicyComponentServiceType;
  public defaultSettingsList: SidenavSetting[] = [
    // notifications
    // { showWarn: true, ...NOTIFICATIONS },
    NOTIFICATIONS,
    SMTP_PROVIDER,
    SMS_PROVIDER,
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
    LANGUAGES,
    OIDC,
    SECRETS,
    SECURITY,
  ];

  public settingsList: Observable<SidenavSetting[]> = of([]);

  private destroy$: Subject<void> = new Subject();
  constructor(
    breadcrumbService: BreadcrumbService,
    activatedRoute: ActivatedRoute,
    public authService: GrpcAuthService,
  ) {
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
    this.settingsList = this.authService.isAllowedMapper(
      this.defaultSettingsList,
      (setting) => setting.requiredRoles.admin || [],
    );
  }

  ngOnDestroy(): void {
    this.destroy$.next();
    this.destroy$.complete();
  }
}
