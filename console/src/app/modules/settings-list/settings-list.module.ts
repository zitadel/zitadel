import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { TranslateModule } from '@ngx-translate/core';
import { HasRolePipeModule } from 'src/app/pipes/has-role-pipe/has-role-pipe.module';

import { CardModule } from '../card/card.module';
import { DomainPolicyModule } from '../policies/domain-policy/domain-policy.module';
import { GeneralSettingsModule } from '../policies/general-settings/general-settings.module';
import { IdpSettingsModule } from '../policies/idp-settings/idp-settings.module';
import { LoginPolicyModule } from '../policies/login-policy/login-policy.module';
import { LoginTextsPolicyModule } from '../policies/login-texts/login-texts.module';
import { MessageTextsPolicyModule } from '../policies/message-texts/message-texts.module';
import { NotificationSettingsModule } from '../policies/notification-settings/notification-settings.module';
import { OIDCConfigurationModule } from '../policies/oidc-configuration/oidc-configuration.module';
import { PasswordComplexityPolicyModule } from '../policies/password-complexity-policy/password-complexity-policy.module';
import { PasswordLockoutPolicyModule } from '../policies/password-lockout-policy/password-lockout-policy.module';
import { PrivacyPolicyModule } from '../policies/privacy-policy/privacy-policy.module';
import { PrivateLabelingPolicyModule } from '../policies/private-labeling-policy/private-labeling-policy.module';
import { SecretGeneratorModule } from '../policies/secret-generator/secret-generator.module';
import { SidenavModule } from '../sidenav/sidenav.module';
import { SettingsListComponent } from './settings-list.component';

@NgModule({
  declarations: [SettingsListComponent],
  imports: [
    CommonModule,
    FormsModule,
    SidenavModule,
    LoginPolicyModule,
    CardModule,
    PasswordComplexityPolicyModule,
    PasswordLockoutPolicyModule,
    PrivateLabelingPolicyModule,
    GeneralSettingsModule,
    IdpSettingsModule,
    PrivacyPolicyModule,
    MessageTextsPolicyModule,
    LoginTextsPolicyModule,
    DomainPolicyModule,
    TranslateModule,
    HasRolePipeModule,
    NotificationSettingsModule,
    OIDCConfigurationModule,
    SecretGeneratorModule,
  ],
  exports: [SettingsListComponent],
})
export class SettingsListModule {}
