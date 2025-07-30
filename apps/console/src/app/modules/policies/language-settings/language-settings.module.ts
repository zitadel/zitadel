import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';
import { FormsModule, ReactiveFormsModule } from '@angular/forms';
import { MatButtonModule } from '@angular/material/button';
import { MatProgressSpinnerModule } from '@angular/material/progress-spinner';
import { MatSelectModule } from '@angular/material/select';
import { TranslateModule } from '@ngx-translate/core';
import { HasRolePipeModule } from 'src/app/pipes/has-role-pipe/has-role-pipe.module';

import { CardModule } from '../../card/card.module';
import { FormFieldModule } from '../../form-field/form-field.module';
import { LanguageSettingsComponent } from './language-settings.component';
import { MatListModule } from '@angular/material/list';
import { MatFormFieldModule } from '@angular/material/form-field';
import { DragDropModule } from '@angular/cdk/drag-drop';
import { MatRadioModule } from '@angular/material/radio';
import { MatTooltipModule } from '@angular/material/tooltip';
import { TableActionsModule } from '../../table-actions/table-actions.module';
import { MatMenuModule } from '@angular/material/menu';
import { MatIconModule } from '@angular/material/icon';

@NgModule({
  declarations: [LanguageSettingsComponent],
  imports: [
    CommonModule,
    CardModule,
    FormsModule,
    ReactiveFormsModule,
    MatFormFieldModule,
    MatButtonModule,
    MatSelectModule,
    FormFieldModule,
    MatProgressSpinnerModule,
    MatSelectModule,
    HasRolePipeModule,
    TranslateModule,
    MatListModule,
    DragDropModule,
    MatRadioModule,
    MatTooltipModule,
    MatMenuModule,
    MatIconModule,
  ],
  exports: [LanguageSettingsComponent],
})
export class LanguageSettingsModule {}
