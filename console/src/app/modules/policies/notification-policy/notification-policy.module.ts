import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { MatIconModule } from '@angular/material/icon';
import { MatLegacyButtonModule as MatButtonModule } from '@angular/material/legacy-button';
import { MatLegacyCheckboxModule as MatCheckboxModule } from '@angular/material/legacy-checkbox';
import { MatLegacyDialogModule as MatDialogModule } from '@angular/material/legacy-dialog';
import { MatLegacyProgressSpinnerModule as MatProgressSpinnerModule } from '@angular/material/legacy-progress-spinner';
import { MatLegacyTooltipModule as MatTooltipModule } from '@angular/material/legacy-tooltip';
import { TranslateModule } from '@ngx-translate/core';
import { HasRoleModule } from 'src/app/directives/has-role/has-role.module';
import { DetailLayoutModule } from 'src/app/modules/detail-layout/detail-layout.module';
import { InputModule } from 'src/app/modules/input/input.module';
import { HasRolePipeModule } from 'src/app/pipes/has-role-pipe/has-role-pipe.module';

import { CardModule } from '../../card/card.module';
import { InfoSectionModule } from '../../info-section/info-section.module';
import { WarnDialogModule } from '../../warn-dialog/warn-dialog.module';
import { NotificationPolicyComponent } from './notification-policy.component';

@NgModule({
  declarations: [NotificationPolicyComponent],
  imports: [
    CommonModule,
    FormsModule,
    InputModule,
    MatButtonModule,
    MatIconModule,
    HasRoleModule,
    MatDialogModule,
    MatTooltipModule,
    MatCheckboxModule,
    HasRolePipeModule,
    TranslateModule,
    WarnDialogModule,
    DetailLayoutModule,
    CardModule,
    MatProgressSpinnerModule,
    InfoSectionModule,
  ],
  exports: [NotificationPolicyComponent],
})
export class NotificationPolicyModule {}
