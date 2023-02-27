import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';
import { FormsModule, ReactiveFormsModule } from '@angular/forms';
import { MatIconModule } from '@angular/material/icon';
import { MatLegacyButtonModule as MatButtonModule } from '@angular/material/legacy-button';
import { MatLegacyProgressSpinnerModule as MatProgressSpinnerModule } from '@angular/material/legacy-progress-spinner';
import { MatLegacySelectModule as MatSelectModule } from '@angular/material/legacy-select';
import { MatLegacySlideToggleModule as MatSlideToggleModule } from '@angular/material/legacy-slide-toggle';
import { TranslateModule } from '@ngx-translate/core';
import { HasRolePipeModule } from 'src/app/pipes/has-role-pipe/has-role-pipe.module';

import { CardModule } from '../../card/card.module';
import { FormFieldModule } from '../../form-field/form-field.module';
import { InputModule } from '../../input/input.module';
import { DialogAddSecretGeneratorComponent } from './dialog-add-secret-generator/dialog-add-secret-generator.component';
import { SecretGeneratorComponent } from './secret-generator.component';

@NgModule({
  declarations: [SecretGeneratorComponent, DialogAddSecretGeneratorComponent],
  imports: [
    CommonModule,
    MatIconModule,
    CardModule,
    FormsModule,
    HasRolePipeModule,
    MatButtonModule,
    FormFieldModule,
    ReactiveFormsModule,
    MatSlideToggleModule,
    InputModule,
    MatProgressSpinnerModule,
    MatSelectModule,
    TranslateModule,
  ],
  exports: [SecretGeneratorComponent, DialogAddSecretGeneratorComponent],
})
export class SecretGeneratorModule {}
