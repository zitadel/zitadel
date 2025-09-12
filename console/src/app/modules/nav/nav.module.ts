import { OverlayModule } from '@angular/cdk/overlay';
import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';
import { FormsModule, ReactiveFormsModule } from '@angular/forms';
import { MatButtonModule } from '@angular/material/button';
import { MatIconModule } from '@angular/material/icon';
import { MatMenuModule } from '@angular/material/menu';
import { MatProgressSpinnerModule } from '@angular/material/progress-spinner';
import { MatTooltipModule } from '@angular/material/tooltip';
import { RouterModule } from '@angular/router';
import { TranslateModule } from '@ngx-translate/core';
import { HasRoleModule } from 'src/app/directives/has-role/has-role.module';
import { HasRolePipeModule } from 'src/app/pipes/has-role-pipe/has-role-pipe.module';

import { MatProgressBarModule } from '@angular/material/progress-bar';
import OnboardingCardModule from '../onboarding-card/onboarding-card.module';
import { NavComponent } from './nav.component';

@NgModule({
  declarations: [NavComponent],
  imports: [
    CommonModule,
    OnboardingCardModule,
    FormsModule,
    ReactiveFormsModule,
    TranslateModule,
    MatIconModule,
    RouterModule,
    MatTooltipModule,
    HasRolePipeModule,
    MatProgressBarModule,
    HasRoleModule,
    MatMenuModule,
    MatButtonModule,
    MatProgressSpinnerModule,
    OverlayModule,
  ],
  exports: [NavComponent],
})
export class NavModule {}
