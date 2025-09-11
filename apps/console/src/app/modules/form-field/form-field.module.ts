import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';
import { MatRippleModule } from '@angular/material/core';
import { TranslateModule } from '@ngx-translate/core';
import { LabelModule } from 'src/app/modules/label/label.module';

import { LabelComponent } from '../label/label.component';
import { CnslErrorDirective } from './error/error.directive';
import { CnslFormFieldComponent } from './field/form-field.component';

@NgModule({
  declarations: [CnslFormFieldComponent, CnslErrorDirective],
  imports: [CommonModule, MatRippleModule, LabelModule, TranslateModule],
  exports: [CnslFormFieldComponent, LabelComponent, CnslErrorDirective],
})
export class FormFieldModule {}
