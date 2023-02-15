import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';
import { MatLegacyButtonModule as MatButtonModule } from '@angular/material/legacy-button';
import { MatRippleModule } from '@angular/material/core';
import { MatIconModule } from '@angular/material/icon';
import { MatLegacyProgressSpinnerModule as MatProgressSpinnerModule } from '@angular/material/legacy-progress-spinner';
import { MatLegacyTooltipModule as MatTooltipModule } from '@angular/material/legacy-tooltip';
import { TranslateModule } from '@ngx-translate/core';
import { HasRoleModule } from 'src/app/directives/has-role/has-role.module';
import { ShortcutsModule } from 'src/app/modules/shortcuts/shortcuts.module';

import { OnboardingCardComponent } from './onboarding-card.component';
import { RouterModule } from '@angular/router';

@NgModule({
  declarations: [OnboardingCardComponent],
  imports: [
    CommonModule,
    MatIconModule,
    HasRoleModule,
    MatButtonModule,
    TranslateModule,
    MatTooltipModule,
    RouterModule,
    MatProgressSpinnerModule,
    ShortcutsModule,
    MatRippleModule,
  ],
  exports: [OnboardingCardComponent],
})
export default class OnboardingCardModule {}
