import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';
import { MatRippleModule } from '@angular/material/core';
import { LabelModule } from 'src/app/modules/label/label.module';

import { LabelComponent } from '../label/label.component';
import { CnslErrorDirective } from './error/error.directive';
import { CnslFormFieldComponent } from './field/form-field.component';
import { I18nErrorsComponent } from './i18n-errors/i18n-errors.component';

@NgModule({
  declarations: [CnslFormFieldComponent, CnslErrorDirective],
  imports: [CommonModule, MatRippleModule, LabelModule],
  exports: [CnslFormFieldComponent, LabelComponent, CnslErrorDirective, I18nErrorsComponent],
})
export class FormFieldModule {}
