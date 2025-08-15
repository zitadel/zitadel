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
import { firstValueFrom, Observable, ReplaySubject, shareReplay, switchMap } from 'rxjs';
import { filter, map, startWith } from 'rxjs/operators';
import { LoginV2FeatureToggleComponent } from '../feature-toggle/login-v2-feature-toggle/login-v2-feature-toggle.component';

// to add a new feature, add the key here and in the FEATURE_KEYS array
const FEATURE_KEYS = [
  'consoleUseV2UserApi',
  'debugOidcParentError',
  'enableBackChannelLogout',
  // 'improvedPerformance',
  'loginDefaultOrg',
  'oidcSingleV1SessionTermination',
  'oidcTokenExchange',
  'permissionCheckV2',
  'userSchema',
] as const;

export type ToggleState = { source: Source; enabled: boolean };
export type ToggleStates = {
  [key in (typeof FEATURE_KEYS)[number]]: ToggleState;
} & {
  loginV2: ToggleState & { baseUri: string };
};

export type ToggleStateKeys = keyof ToggleStates;

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
    LoginV2FeatureToggleComponent,
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
  protected readonly FEATURE_KEYS = FEATURE_KEYS;

  constructor(
    private readonly featureService: NewFeatureService,
    private readonly breadcrumbService: BreadcrumbService,
    private readonly toast: ToastService,
  ) {
    const breadcrumbs = [
      new Breadcrumb({
        type: BreadcrumbType.INSTANCE,
        name: 'Instance',
        routerLink: ['/instance'],
      }),
    ];
    this.breadcrumbService.setBreadcrumb(breadcrumbs);

    this.toggleStates$ = this.getToggleStates().pipe(shareReplay({ refCount: true, bufferSize: 1 }));
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
    return FEATURE_KEYS.reduce(
      (acc, key) => {
        const feature = featureData[key];
        acc[key] = {
          source: feature?.source ?? Source.SYSTEM,
          enabled: !!feature?.enabled,
        };
        return acc;
      },
      {
        // to add special feature flags they have to be mapped here
        loginV2: {
          source: featureData.loginV2?.source ?? Source.SYSTEM,
          enabled: !!featureData.loginV2?.required,
          baseUri: featureData.loginV2?.baseUri ?? '',
        },
      } as ToggleStates,
    );
  }

  public async saveFeatures<TKey extends ToggleStateKeys, TValue extends ToggleStates[TKey]>(key: TKey, value: TValue) {
    const toggleStates = { ...(await firstValueFrom(this.toggleStates$)), [key]: value };

    const req = FEATURE_KEYS.reduce<MessageInitShape<typeof SetInstanceFeaturesRequestSchema>>((acc, key) => {
      acc[key] = toggleStates[key].enabled;
      return acc;
    }, {});

    // to save special flags they have to be handled here
    req['loginV2'] = {
      required: toggleStates.loginV2.enabled,
      baseUri: toggleStates.loginV2.baseUri,
    };

    try {
      await this.featureService.setInstanceFeatures(req);

      // needed because of eventual consistency
      await new Promise((res) => setTimeout(res, 1000));
      this.refresh$.next(true);

      this.toast.showInfo('POLICY.TOAST.SET', true);
    } catch (error) {
      this.toast.showError(error);
    }
  }

  public async resetFeatures() {
    try {
      await this.featureService.resetInstanceFeatures();

      // needed because of eventual consistency
      await new Promise((res) => setTimeout(res, 1000));
      this.refresh$.next(true);

      this.toast.showInfo('POLICY.TOAST.RESETSUCCESS', true);
    } catch (error) {
      this.toast.showError(error);
    }
  }
}
