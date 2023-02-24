import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';
import { MatRippleModule } from '@angular/material/core';
import { MatIconModule } from '@angular/material/icon';
import { MatLegacyProgressSpinnerModule as MatProgressSpinnerModule } from '@angular/material/legacy-progress-spinner';
import { MatLegacyTooltipModule as MatTooltipModule } from '@angular/material/legacy-tooltip';
import { TranslateModule } from '@ngx-translate/core';
import { ShortcutsModule } from 'src/app/modules/shortcuts/shortcuts.module';

import { MatLegacyProgressBarModule } from '@angular/material/legacy-progress-bar';
import { RouterModule } from '@angular/router';
import { EventPipeModule } from 'src/app/pipes/event-pipe/event-pipe.module';
import { OnboardingComponent } from './onboarding.component';

@NgModule({
  declarations: [OnboardingComponent],
  imports: [
    CommonModule,
    MatIconModule,
    TranslateModule,
    MatTooltipModule,
    ShortcutsModule,
    MatRippleModule,
    RouterModule,
    MatProgressSpinnerModule,
    MatLegacyProgressBarModule,
    EventPipeModule,
  ],
  exports: [OnboardingComponent],
})
export default class OnboardingModule {}
