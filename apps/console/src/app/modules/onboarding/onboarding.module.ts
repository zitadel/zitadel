import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';
import { MatRippleModule } from '@angular/material/core';
import { MatIconModule } from '@angular/material/icon';
import { MatProgressSpinnerModule } from '@angular/material/progress-spinner';
import { MatTooltipModule } from '@angular/material/tooltip';
import { TranslateModule } from '@ngx-translate/core';
import { ShortcutsModule } from 'src/app/modules/shortcuts/shortcuts.module';

import { MatProgressBarModule } from '@angular/material/progress-bar';
import { RouterModule } from '@angular/router';
import { MilestonePipeModule } from 'src/app/pipes/milestone-pipe/milestone-pipe.module';
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
    MatProgressBarModule,
    MilestonePipeModule,
  ],
  exports: [OnboardingComponent],
})
export default class OnboardingModule {}
