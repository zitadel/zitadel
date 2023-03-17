import { TextFieldModule } from '@angular/cdk/text-field';
import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';
import { FormsModule, ReactiveFormsModule } from '@angular/forms';
import { MatIconModule } from '@angular/material/icon';
import { MatLegacyButtonModule as MatButtonModule } from '@angular/material/legacy-button';
import { MatLegacyProgressSpinnerModule as MatProgressSpinnerModule } from '@angular/material/legacy-progress-spinner';
import { MatLegacySelectModule as MatSelectModule } from '@angular/material/legacy-select';
import { MatLegacyTooltipModule as MatTooltipModule } from '@angular/material/legacy-tooltip';
import { TranslateModule } from '@ngx-translate/core';
import { HasRoleModule } from 'src/app/directives/has-role/has-role.module';
import { DetailLayoutModule } from 'src/app/modules/detail-layout/detail-layout.module';
import { InputModule } from 'src/app/modules/input/input.module';
import { HasRolePipeModule } from 'src/app/pipes/has-role-pipe/has-role-pipe.module';

import { CardModule } from '../../card/card.module';
import { EditTextModule } from '../../edit-text/edit-text.module';
import { FormFieldModule } from '../../form-field/form-field.module';
import { InfoSectionModule } from '../../info-section/info-section.module';
import { MessageTextsRoutingModule } from './message-texts-routing.module';
import { MessageTextsComponent } from './message-texts.component';

@NgModule({
  declarations: [MessageTextsComponent],
  imports: [
    MessageTextsRoutingModule,
    CommonModule,
    InfoSectionModule,
    ReactiveFormsModule,
    FormsModule,
    InputModule,
    FormFieldModule,
    EditTextModule,
    MatButtonModule,
    MatIconModule,
    HasRoleModule,
    HasRolePipeModule,
    MatTooltipModule,
    TranslateModule,
    CardModule,
    MatTooltipModule,
    MatSelectModule,
    DetailLayoutModule,
    MatProgressSpinnerModule,
    TextFieldModule,
  ],
  exports: [MessageTextsComponent],
})
export class MessageTextsPolicyModule {}
