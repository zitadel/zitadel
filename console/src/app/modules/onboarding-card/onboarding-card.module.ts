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

import { RouterModule } from '@angular/router';
import { EventPipeModule } from 'src/app/pipes/event-pipe/event-pipe.module';
import { LocalizedDatePipeModule } from 'src/app/pipes/localized-date-pipe/localized-date-pipe.module';
import { TimestampToDatePipeModule } from 'src/app/pipes/timestamp-to-date-pipe/timestamp-to-date-pipe.module';
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
