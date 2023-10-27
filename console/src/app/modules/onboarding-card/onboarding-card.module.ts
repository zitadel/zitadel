import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';
import { MatIconModule } from '@angular/material/icon';
import { MatProgressSpinnerModule } from '@angular/material/progress-spinner';
import { MatTooltipModule } from '@angular/material/tooltip';
import { TranslateModule } from '@ngx-translate/core';

import { RouterModule } from '@angular/router';
import { MilestonePipeModule } from 'src/app/pipes/milestone-pipe/milestone-pipe.module';
import { OnboardingCardComponent } from './onboarding-card.component';

@NgModule({
  declarations: [OnboardingCardComponent],
  imports: [
    CommonModule,
    MatIconModule,
    TranslateModule,
    RouterModule,
    MatProgressSpinnerModule,
    MilestonePipeModule,
    MatTooltipModule,
  ],
  exports: [OnboardingCardComponent],
})
export default class OnboardingCardModule {}
