import { Component, OnInit } from '@angular/core';
import { ActivatedRoute, Params } from '@angular/router';
import { Observable, of, take } from 'rxjs';
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
  NOTIFICATIONS,
  PRIVACYPOLICY,
  VERIFIED_DOMAINS,
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
    NOTIFICATIONS,
    VERIFIED_DOMAINS,
    DOMAIN,
    BRANDING,
    MESSAGETEXTS,
    LOGINTEXTS,
    PRIVACYPOLICY,
  ];

  public settingsList: Observable<Array<SidenavSetting>> = of([]);

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
    this.settingsList = this.authService
      .isAllowedMapper(this.defaultSettingsList, (setting) => setting.requiredRoles.mgmt || [])
      .pipe(take(1));
  }
}
