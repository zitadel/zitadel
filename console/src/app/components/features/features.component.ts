import { CommonModule } from '@angular/common';
import { Component } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { MatButtonModule } from '@angular/material/button';
import { MatButtonToggleModule } from '@angular/material/button-toggle';
import { MatCheckboxModule } from '@angular/material/checkbox';
import { MatIconModule } from '@angular/material/icon';
import { MatTooltipModule } from '@angular/material/tooltip';
import { TranslateModule } from '@ngx-translate/core';
import { HasRoleModule } from 'src/app/directives/has-role/has-role.module';
import { CardModule } from 'src/app/modules/card/card.module';
import { InfoSectionModule } from 'src/app/modules/info-section/info-section.module';
import { HasRolePipeModule } from 'src/app/pipes/has-role-pipe/has-role-pipe.module';
import { Breadcrumb, BreadcrumbService, BreadcrumbType } from 'src/app/services/breadcrumb.service';
import { ToastService } from 'src/app/services/toast.service';
import { FeatureToggleComponent } from '../feature-toggle/feature-toggle.component';
import { NewFeatureService } from 'src/app/services/new-feature.service';
import {
  GetInstanceFeaturesResponse,
  SetInstanceFeaturesRequestSchema,
} from '@zitadel/proto/zitadel/feature/v2/instance_pb';
import { Source } from '@zitadel/proto/zitadel/feature/v2/feature_pb';
import { MessageInitShape } from '@bufbuild/protobuf';
import { Observable, ReplaySubject, switchMap } from 'rxjs';
import { filter, map, startWith } from 'rxjs/operators';

// TODO: to add a new feature, add the key here and in the FEATURE_KEYS array
const FEATURE_KEYS = [
  'actions',
  'consoleUseV2UserApi',
  'debugOidcParentError',
  'disableUserTokenEvent',
  'enableBackChannelLogout',
  // 'improvedPerformance',
  'loginDefaultOrg',
  // 'loginV2',
  'oidcLegacyIntrospection',
  'oidcSingleV1SessionTermination',
  'oidcTokenExchange',
  'oidcTriggerIntrospectionProjections',
  'permissionCheckV2',
  'userSchema',
  // 'webKey',
] as const;

export type FeatureState = { source: Source; enabled: boolean };
export type ToggleStateKeys = (typeof FEATURE_KEYS)[number];

export type ToggleStates = {
  [key in ToggleStateKeys]: FeatureState;
};

@Component({
  imports: [
    CommonModule,
    FormsModule,
    MatButtonToggleModule,
    HasRolePipeModule,
    MatIconModule,
    CardModule,
    TranslateModule,
    MatButtonModule,
    MatCheckboxModule,
    InfoSectionModule,
    MatTooltipModule,
    HasRoleModule,
    FeatureToggleComponent,
  ],
  standalone: true,
  selector: 'cnsl-features',
  templateUrl: './features.component.html',
  styleUrls: ['./features.component.scss'],
})
export class FeaturesComponent {
  private readonly refresh$ = new ReplaySubject<true>(1);
  protected readonly toggleStates$: Observable<ToggleStates>;
  protected readonly Source = Source;

  constructor(
    private featureService: NewFeatureService,
    private breadcrumbService: BreadcrumbService,
    private toast: ToastService,
  ) {
    const breadcrumbs = [
      new Breadcrumb({
        type: BreadcrumbType.INSTANCE,
        name: 'Instance',
        routerLink: ['/instance'],
      }),
    ];
    this.breadcrumbService.setBreadcrumb(breadcrumbs);

    this.toggleStates$ = this.getToggleStates();
  }

  public async validateAndSave(toggleStates: ToggleStates) {
    const req = FEATURE_KEYS.reduce<MessageInitShape<typeof SetInstanceFeaturesRequestSchema>>((acc, key) => {
      acc[key] = toggleStates[key].enabled;
      return acc;
    }, {});

    try {
      await this.featureService.setInstanceFeatures(req);
      this.toast.showInfo('POLICY.TOAST.SET', true);
    } catch (error) {
      this.toast.showError(error);
    }
  }

  private getToggleStates() {
    return this.refresh$.pipe(
      startWith(true),
      switchMap(async () => {
        try {
          return await this.featureService.getInstanceFeatures();
        } catch (error) {
          this.toast.showError(error);
          return undefined;
        }
      }),
      filter(Boolean),
      map((res) => this.createToggleStates(res)),
    );
  }

  private createToggleStates(featureData: GetInstanceFeaturesResponse): ToggleStates {
    return FEATURE_KEYS.reduce((acc, key) => {
      const feature = featureData[key];
      acc[key] = {
        source: feature?.source ?? Source.SYSTEM,
        enabled: !!feature?.enabled,
      };
      return acc;
    }, {} as ToggleStates);
  }

  public async resetSettings() {
    try {
      await this.featureService.resetInstanceFeatures();
      this.toast.showInfo('POLICY.TOAST.RESETSUCCESS', true);

      await new Promise((res) => setTimeout(res, 1000));
      this.refresh$.next(true);
    } catch (error) {
      this.toast.showError(error);
    }
  }

  public get toggleStateKeys() {
    return FEATURE_KEYS;
  }
}
