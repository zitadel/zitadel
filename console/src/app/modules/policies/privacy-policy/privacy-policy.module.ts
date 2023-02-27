import { TextFieldModule } from '@angular/cdk/text-field';
import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';
import { FormsModule, ReactiveFormsModule } from '@angular/forms';
import { MatIconModule } from '@angular/material/icon';
import { MatLegacyButtonModule as MatButtonModule } from '@angular/material/legacy-button';
import { MatLegacyDialogModule as MatDialogModule } from '@angular/material/legacy-dialog';
import { MatLegacyProgressSpinnerModule as MatProgressSpinnerModule } from '@angular/material/legacy-progress-spinner';
import { MatLegacySelectModule as MatSelectModule } from '@angular/material/legacy-select';
import { MatLegacyTooltipModule as MatTooltipModule } from '@angular/material/legacy-tooltip';
import { TranslateModule } from '@ngx-translate/core';
import { CopyToClipboardModule } from 'src/app/directives/copy-to-clipboard/copy-to-clipboard.module';

import { HasRoleModule } from '../../../directives/has-role/has-role.module';
import { DetailLayoutModule } from '../../../modules/detail-layout/detail-layout.module';
import { InputModule } from '../../../modules/input/input.module';
import { HasRolePipeModule } from '../../../pipes/has-role-pipe/has-role-pipe.module';
import { CardModule } from '../../card/card.module';
import { FormFieldModule } from '../../form-field/form-field.module';
import { InfoSectionModule } from '../../info-section/info-section.module';
import { WarnDialogModule } from '../../warn-dialog/warn-dialog.module';
import { PrivacyPolicyRoutingModule } from './privacy-policy-routing.module';
import { PrivacyPolicyComponent } from './privacy-policy.component';

@NgModule({
  declarations: [PrivacyPolicyComponent],
  imports: [
    PrivacyPolicyRoutingModule,
    MatSelectModule,
    CommonModule,
    ReactiveFormsModule,
    FormsModule,
    InputModule,
    FormFieldModule,
    CopyToClipboardModule,
    MatButtonModule,
    MatIconModule,
    HasRoleModule,
    HasRolePipeModule,
    MatTooltipModule,
    TranslateModule,
    CardModule,
    MatTooltipModule,
    DetailLayoutModule,
    MatProgressSpinnerModule,
    TextFieldModule,
    MatDialogModule,
    WarnDialogModule,
    InfoSectionModule,
  ],
  exports: [PrivacyPolicyComponent],
})
export class PrivacyPolicyModule {}
