import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';
import { MatIconModule } from '@angular/material/icon';
import { MatLegacyProgressSpinnerModule as MatProgressSpinnerModule } from '@angular/material/legacy-progress-spinner';
import { MatLegacyTooltipModule as MatTooltipModule } from '@angular/material/legacy-tooltip';
import { TranslateModule } from '@ngx-translate/core';

import { RouterModule } from '@angular/router';
import { EventPipeModule } from 'src/app/pipes/event-pipe/event-pipe.module';
import { OnboardingCardComponent } from './onboarding-card.component';

@NgModule({
  declarations: [OnboardingCardComponent],
  imports: [
    CommonModule,
    MatIconModule,
    TranslateModule,
    RouterModule,
    MatProgressSpinnerModule,
    EventPipeModule,
    MatTooltipModule,
  ],
  exports: [OnboardingCardComponent],
})
export default class OnboardingCardModule {}
