import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { MatButtonModule } from '@angular/material/button';
import { TranslateModule } from '@ngx-translate/core';

import { InputModule } from '../input/input.module';
import { WarnDialogComponent } from './warn-dialog.component';

@NgModule({
  declarations: [WarnDialogComponent],
  imports: [CommonModule, FormsModule, TranslateModule, MatButtonModule, InputModule],
})
export class WarnDialogModule {}
