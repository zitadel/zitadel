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
import { Source } from 'src/app/proto/generated/zitadel/feature/v2beta/feature_pb';
import { Breadcrumb, BreadcrumbService, BreadcrumbType } from 'src/app/services/breadcrumb.service';
import { ToastService } from 'src/app/services/toast.service';
import { FeatureToggleComponent } from '../feature-toggle/feature-toggle.component';
import { FeatureFlag } from 'src/app/proto/generated/zitadel/feature/v2/feature_pb';
import { NewFeatureService } from 'src/app/services/new-feature.service';
import {
  GetInstanceFeaturesResponse,
  SetInstanceFeaturesRequest,
  SetInstanceFeaturesRequestSchema,
} from '@zitadel/proto/zitadel/feature/v2/instance_pb';
//@ts-ignore
import { create } from '@zitadel/client';

export enum ToggleState {
  ENABLED = 'ENABLED',
  DISABLED = 'DISABLED',
}

// TODO: to add a new feature, add the key here and in the FEATURE_KEYS array
const FEATURE_KEYS: ToggleStateKeys[] = [
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
  'webKey',
];

type FeatureState = { source: Source; state: ToggleState };
export type ToggleStateKeys = Exclude<keyof GetInstanceFeaturesResponse, 'details' | '$typeName' | '$unknown'>;

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
  protected featureData: GetInstanceFeaturesResponse | undefined;

  protected toggleStates: ToggleStates | undefined;
  protected Source: any = Source;
  protected ToggleState: any = ToggleState;

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

    this.getFeatures(true);
  }

  public validateAndSave() {
    const req = create(SetInstanceFeaturesRequestSchema, {
      actions: this.toggleStates?.actions?.state === ToggleState.ENABLED,
      consoleUseV2UserApi: this.toggleStates?.consoleUseV2UserApi?.state === ToggleState.ENABLED,
      debugOidcParentError: this.toggleStates?.debugOidcParentError?.state === ToggleState.ENABLED,
      disableUserTokenEvent: this.toggleStates?.disableUserTokenEvent?.state === ToggleState.ENABLED,
      enableBackChannelLogout: this.toggleStates?.enableBackChannelLogout?.state === ToggleState.ENABLED,
      loginDefaultOrg: this.toggleStates?.loginDefaultOrg?.state === ToggleState.ENABLED,
      oidcLegacyIntrospection: this.toggleStates?.oidcLegacyIntrospection?.state === ToggleState.ENABLED,
      oidcSingleV1SessionTermination: this.toggleStates?.oidcSingleV1SessionTermination?.state === ToggleState.ENABLED,
      oidcTokenExchange: this.toggleStates?.oidcTokenExchange?.state === ToggleState.ENABLED,
      oidcTriggerIntrospectionProjections:
        this.toggleStates?.oidcTriggerIntrospectionProjections?.state === ToggleState.ENABLED,
      permissionCheckV2: this.toggleStates?.permissionCheckV2?.state === ToggleState.ENABLED,
      userSchema: this.toggleStates?.userSchema?.state === ToggleState.ENABLED,
      webKey: this.toggleStates?.webKey?.state === ToggleState.ENABLED,
    });

    this.featureService
      .setInstanceFeatures(req)
      .then(() => {
        this.toast.showInfo('POLICY.TOAST.SET', true);
      })
      .catch((error) => {
        this.toast.showError(error);
      });
  }

  private getFeatures(inheritance: boolean) {
    this.featureService.getInstanceFeatures().then((instanceFeaturesResponse) => {
      this.featureData = instanceFeaturesResponse;

      console.log(this.featureData);
      this.toggleStates = this.createToggleStates(this.featureData);
    });
  }

  private createToggleStates(featureData: GetInstanceFeaturesResponse): ToggleStates {
    const toggleStates: Partial<ToggleStates> = {};

    FEATURE_KEYS.forEach((key) => {
      // TODO: Fix this type cast as not all keys are present as FeatureFlag
      const feature = featureData[key] as unknown as FeatureFlag.AsObject;
      toggleStates[key] = {
        source: feature?.source || Source.SOURCE_SYSTEM,
        state: !!feature?.enabled ? ToggleState.ENABLED : ToggleState.DISABLED,
      };
    });

    return toggleStates as ToggleStates;
  }

  public resetSettings(): void {
    this.featureService
      .resetInstanceFeatures()
      .then(() => {
        this.toast.showInfo('POLICY.TOAST.RESETSUCCESS', true);
        setTimeout(() => {
          this.getFeatures(true);
        }, 1000);
      })
      .catch((error) => {
        this.toast.showError(error);
      });
  }

  public get toggleStateKeys() {
    return Object.keys(this.toggleStates ?? {});
  }
}
