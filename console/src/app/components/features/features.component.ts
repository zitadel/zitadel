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
import { FeatureService } from 'src/app/services/feature.service';
import { ToastService } from 'src/app/services/toast.service';
import {
  GetInstanceFeaturesResponse,
  SetInstanceFeaturesRequest,
} from 'src/app/proto/generated/zitadel/feature/v2/instance_pb';

enum ToggleState {
  ENABLED = 'ENABLED',
  DISABLED = 'DISABLED',
  INHERITED = 'INHERITED',
}

type FeatureState = { source: Source; state: ToggleState };
type ToggleStates = {
  loginDefaultOrg?: FeatureState;
  oidcTriggerIntrospectionProjections?: FeatureState;
  oidcLegacyIntrospection?: FeatureState;
  userSchema?: FeatureState;
  oidcTokenExchange?: FeatureState;
  actions?: FeatureState;
  oidcSingleV1SessionTermination?: FeatureState;
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
  ],
  standalone: true,
  selector: 'cnsl-features',
  templateUrl: './features.component.html',
  styleUrls: ['./features.component.scss'],
})
export class FeaturesComponent {
  protected featureData: GetInstanceFeaturesResponse.AsObject | undefined;

  protected toggleStates: ToggleStates | undefined;
  protected Source: any = Source;
  protected ToggleState: any = ToggleState;

  constructor(
    private featureService: FeatureService,
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
    this.featureService.resetInstanceFeatures().then(() => {
      const req = new SetInstanceFeaturesRequest();
      let changed = false;

      console.log(this.toggleStates);

      if (this.toggleStates?.loginDefaultOrg?.state !== ToggleState.INHERITED) {
        req.setLoginDefaultOrg(this.toggleStates?.loginDefaultOrg?.state === ToggleState.ENABLED);
        changed = true;
      }
      if (this.toggleStates?.oidcTriggerIntrospectionProjections?.state !== ToggleState.INHERITED) {
        req.setOidcTriggerIntrospectionProjections(
          this.toggleStates?.oidcTriggerIntrospectionProjections?.state === ToggleState.ENABLED,
        );
        changed = true;
      }
      if (this.toggleStates?.oidcLegacyIntrospection?.state !== ToggleState.INHERITED) {
        req.setOidcLegacyIntrospection(this.toggleStates?.oidcLegacyIntrospection?.state === ToggleState.ENABLED);
        changed = true;
      }
      if (this.toggleStates?.userSchema?.state !== ToggleState.INHERITED) {
        req.setUserSchema(this.toggleStates?.userSchema?.state === ToggleState.ENABLED);
        changed = true;
      }
      if (this.toggleStates?.oidcTokenExchange?.state !== ToggleState.INHERITED) {
        req.setOidcTokenExchange(this.toggleStates?.oidcTokenExchange?.state === ToggleState.ENABLED);
        changed = true;
      }
      if (this.toggleStates?.actions?.state !== ToggleState.INHERITED) {
        req.setActions(this.toggleStates?.actions?.state === ToggleState.ENABLED);
        changed = true;
      }
      if (this.toggleStates?.oidcSingleV1SessionTermination?.state !== ToggleState.INHERITED) {
        req.setOidcSingleV1SessionTermination(
          this.toggleStates?.oidcSingleV1SessionTermination?.state === ToggleState.ENABLED,
        );
        changed = true;
      }

      if (changed) {
        this.featureService
          .setInstanceFeatures(req)
          .then(() => {
            this.toast.showInfo('POLICY.TOAST.SET', true);
          })
          .catch((error) => {
            this.toast.showError(error);
          });
      }
    });
  }

  private getFeatures(inheritance: boolean) {
    this.featureService.getInstanceFeatures(inheritance).then((instanceFeaturesResponse) => {
      this.featureData = instanceFeaturesResponse.toObject();
      console.log(this.featureData);

      this.toggleStates = {
        loginDefaultOrg: {
          source: this.featureData.loginDefaultOrg?.source || Source.SOURCE_SYSTEM,
          state:
            this.featureData.loginDefaultOrg?.source === Source.SOURCE_SYSTEM ||
            this.featureData.loginDefaultOrg?.source === Source.SOURCE_UNSPECIFIED
              ? ToggleState.INHERITED
              : !!this.featureData.loginDefaultOrg?.enabled
                ? ToggleState.ENABLED
                : ToggleState.DISABLED,
        },
        oidcTriggerIntrospectionProjections: {
          source: this.featureData.oidcTriggerIntrospectionProjections?.source || Source.SOURCE_SYSTEM,
          state:
            this.featureData.oidcTriggerIntrospectionProjections?.source === Source.SOURCE_SYSTEM ||
            this.featureData.oidcTriggerIntrospectionProjections?.source === Source.SOURCE_UNSPECIFIED
              ? ToggleState.INHERITED
              : !!this.featureData.oidcTriggerIntrospectionProjections?.enabled
                ? ToggleState.ENABLED
                : ToggleState.DISABLED,
        },
        oidcLegacyIntrospection: {
          source: this.featureData.oidcLegacyIntrospection?.source || Source.SOURCE_SYSTEM,
          state:
            this.featureData.oidcLegacyIntrospection?.source === Source.SOURCE_SYSTEM ||
            this.featureData.oidcLegacyIntrospection?.source === Source.SOURCE_UNSPECIFIED
              ? ToggleState.INHERITED
              : !!this.featureData.oidcLegacyIntrospection?.enabled
                ? ToggleState.ENABLED
                : ToggleState.DISABLED,
        },
        userSchema: {
          source: this.featureData.userSchema?.source || Source.SOURCE_SYSTEM,
          state:
            this.featureData.userSchema?.source === Source.SOURCE_SYSTEM ||
            this.featureData.userSchema?.source === Source.SOURCE_UNSPECIFIED
              ? ToggleState.INHERITED
              : !!this.featureData.userSchema?.enabled
                ? ToggleState.ENABLED
                : ToggleState.DISABLED,
        },
        oidcTokenExchange: {
          source: this.featureData.oidcTokenExchange?.source || Source.SOURCE_SYSTEM,
          state:
            this.featureData.oidcTokenExchange?.source === Source.SOURCE_SYSTEM ||
            this.featureData.oidcTokenExchange?.source === Source.SOURCE_UNSPECIFIED
              ? ToggleState.INHERITED
              : !!this.featureData.oidcTokenExchange?.enabled
                ? ToggleState.ENABLED
                : ToggleState.DISABLED,
        },
        actions: {
          source: Source.SOURCE_SYSTEM,
          state:
            this.featureData.actions?.source === Source.SOURCE_SYSTEM ||
            this.featureData.actions?.source === Source.SOURCE_UNSPECIFIED
              ? ToggleState.INHERITED
              : !!this.featureData.actions?.enabled
                ? ToggleState.ENABLED
                : ToggleState.DISABLED,
        },
        oidcSingleV1SessionTermination: {
          source: this.featureData.oidcSingleV1SessionTermination?.source || Source.SOURCE_SYSTEM,
          state:
            this.featureData.oidcSingleV1SessionTermination?.source === Source.SOURCE_SYSTEM ||
            this.featureData.oidcSingleV1SessionTermination?.source === Source.SOURCE_UNSPECIFIED
              ? ToggleState.INHERITED
              : !!this.featureData.oidcSingleV1SessionTermination?.enabled
                ? ToggleState.ENABLED
                : ToggleState.DISABLED,
        },
      };
    });
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
}
