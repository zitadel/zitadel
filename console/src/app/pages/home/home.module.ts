import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';
import { MatButtonModule } from '@angular/material/button';
import { MatRippleModule } from '@angular/material/core';
import { MatIconModule } from '@angular/material/icon';
import { MatProgressSpinnerModule } from '@angular/material/progress-spinner';
import { MatTooltipModule } from '@angular/material/tooltip';
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
