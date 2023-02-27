import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';
import { MatIconModule } from '@angular/material/icon';
import { MatLegacyButtonModule as MatButtonModule } from '@angular/material/legacy-button';
import { MatLegacyProgressBarModule as MatProgressBarModule } from '@angular/material/legacy-progress-bar';
import { RouterModule } from '@angular/router';
import { TranslateModule } from '@ngx-translate/core';

import { AvatarModule } from '../avatar/avatar.module';
import { AccountsCardComponent } from './accounts-card.component';

@NgModule({
  declarations: [AccountsCardComponent],
  imports: [CommonModule, MatIconModule, MatButtonModule, MatProgressBarModule, RouterModule, AvatarModule, TranslateModule],
  exports: [AccountsCardComponent],
})
export class AccountsCardModule {}
