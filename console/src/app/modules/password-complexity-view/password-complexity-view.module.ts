import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { MatProgressSpinnerModule } from '@angular/material/progress-spinner';
import { TranslateModule } from '@ngx-translate/core';

import { PasswordComplexityViewComponent } from './password-complexity-view.component';

@NgModule({
  declarations: [PasswordComplexityViewComponent],
  imports: [CommonModule, MatProgressSpinnerModule, TranslateModule, FormsModule],
  exports: [PasswordComplexityViewComponent],
})
export class PasswordComplexityViewModule {}
