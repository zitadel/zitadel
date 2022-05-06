import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { TranslateModule } from '@ngx-translate/core';
import { HasRolePipeModule } from 'src/app/pipes/has-role-pipe/has-role-pipe.module';

import { CardModule } from '../card/card.module';
import { IdpTableModule } from '../idp-table/idp-table.module';
import { GeneralSettingsModule } from '../policies/general-settings/general-settings.module';
import { LoginPolicyModule } from '../policies/login-policy/login-policy.module';
import { LoginTextsPolicyModule } from '../policies/login-texts/login-texts.module';
import { MessageTextsPolicyModule } from '../policies/message-texts/message-texts.module';
import { OrgIamPolicyModule } from '../policies/org-iam-policy/org-iam-policy.module';
import { PasswordComplexityPolicyModule } from '../policies/password-complexity-policy/password-complexity-policy.module';
import { PasswordLockoutPolicyModule } from '../policies/password-lockout-policy/password-lockout-policy.module';
import { PrivacyPolicyModule } from '../policies/privacy-policy/privacy-policy.module';
import { PrivateLabelingPolicyModule } from '../policies/private-labeling-policy/private-labeling-policy.module';
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
    IdpTableModule,
    PrivacyPolicyModule,
    MessageTextsPolicyModule,
    LoginTextsPolicyModule,
    OrgIamPolicyModule,
    TranslateModule,
    HasRolePipeModule,
  ],
  exports: [SettingsListComponent],
})
export class SettingsListModule {}
