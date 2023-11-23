import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';
import {FormsModule, ReactiveFormsModule} from '@angular/forms';
import { MatButtonModule } from '@angular/material/button';
import { MatProgressSpinnerModule } from '@angular/material/progress-spinner';
import { MatSelectModule } from '@angular/material/select';
import { TranslateModule } from '@ngx-translate/core';
import { HasRolePipeModule } from 'src/app/pipes/has-role-pipe/has-role-pipe.module';

import { CardModule } from '../../card/card.module';
import { FormFieldModule } from '../../form-field/form-field.module';
import { LanguageSettingsComponent } from './language-settings.component';
import { MatListModule} from "@angular/material/list";
import { MatFormFieldModule} from "@angular/material/form-field";
import { DragDropModule} from "@angular/cdk/drag-drop";

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
  ],
  exports: [LanguageSettingsComponent],
})
export class LanguageSettingsModule {}
