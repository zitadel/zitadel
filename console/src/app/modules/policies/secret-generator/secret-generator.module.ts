import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';
import { FormsModule, ReactiveFormsModule } from '@angular/forms';
import { MatIconModule } from '@angular/material/icon';
import { MatButtonModule } from '@angular/material/button';
import { MatProgressSpinnerModule } from '@angular/material/progress-spinner';
import { MatSelectModule } from '@angular/material/select';
import { MatSlideToggleModule } from '@angular/material/slide-toggle';
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
