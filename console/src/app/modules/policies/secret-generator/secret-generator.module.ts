import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';
import { FormsModule, ReactiveFormsModule } from '@angular/forms';
import { MatButtonModule } from '@angular/material/button';
import { MatIconModule } from '@angular/material/icon';
import { MatProgressSpinnerModule } from '@angular/material/progress-spinner';
import { MatSelectModule } from '@angular/material/select';
import { MatSlideToggleModule } from '@angular/material/slide-toggle';
import { TranslateModule } from '@ngx-translate/core';
import { HasRolePipeModule } from 'src/app/pipes/has-role-pipe/has-role-pipe.module';
import { MatDialogModule } from '@angular/material/dialog';
import { CardModule } from '../../card/card.module';
import { FormFieldModule } from '../../form-field/form-field.module';
import { InputModule } from '../../input/input.module';
import { SecretGeneratorComponent } from './secret-generator.component';
import { SecretGeneratorCardComponent } from './secret-generator-card/secret-generator-card.component';

@NgModule({
  declarations: [SecretGeneratorComponent, SecretGeneratorCardComponent],
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
    MatDialogModule,
  ],
  exports: [SecretGeneratorComponent],
})
export class SecretGeneratorModule {}
