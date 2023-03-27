import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';
import { MatRippleModule } from '@angular/material/core';
import { MatIconModule } from '@angular/material/icon';
import { MatLegacyButtonModule as MatButtonModule } from '@angular/material/legacy-button';
import { MatLegacyProgressSpinnerModule as MatProgressSpinnerModule } from '@angular/material/legacy-progress-spinner';
import { MatLegacyTooltipModule as MatTooltipModule } from '@angular/material/legacy-tooltip';
import { TranslateModule } from '@ngx-translate/core';
import { HasRoleModule } from 'src/app/directives/has-role/has-role.module';
import { ShortcutsModule } from 'src/app/modules/shortcuts/shortcuts.module';

import OnboardingModule from 'src/app/modules/onboarding/onboarding.module';
import { HasRolePipeModule } from 'src/app/pipes/has-role-pipe/has-role-pipe.module';
import { HomeRoutingModule } from './home-routing.module';
import { HomeComponent } from './home.component';

@NgModule({
  declarations: [HomeComponent],
  imports: [
    CommonModule,
    MatIconModule,
    HasRoleModule,
    HomeRoutingModule,
    MatButtonModule,
    HasRolePipeModule,
    TranslateModule,
    MatTooltipModule,
    MatProgressSpinnerModule,
    ShortcutsModule,
    OnboardingModule,
    MatRippleModule,
  ],
})
export default class HomeModule {}
