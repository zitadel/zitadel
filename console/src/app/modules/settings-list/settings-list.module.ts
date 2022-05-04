import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';
import { FormsModule } from '@angular/forms';

import { LoginPolicyModule } from '../policies/login-policy/login-policy.module';
import { PasswordComplexityPolicyModule } from '../policies/password-complexity-policy/password-complexity-policy.module';
import { PasswordLockoutPolicyModule } from '../policies/password-lockout-policy/password-lockout-policy.module';
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
    PasswordComplexityPolicyModule,
    PasswordLockoutPolicyModule,
    PrivateLabelingPolicyModule,
  ],
  exports: [SettingsListComponent],
})
export class SettingsListModule {}
