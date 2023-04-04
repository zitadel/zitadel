import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';
import { FormsModule, ReactiveFormsModule } from '@angular/forms';
import { MatIconModule } from '@angular/material/icon';
import { MatLegacyButtonModule } from '@angular/material/legacy-button';
import { MatLegacyTooltipModule } from '@angular/material/legacy-tooltip';
import { TranslateModule } from '@ngx-translate/core';
import { InputModule } from '../input/input.module';
import { StringListComponent } from './string-list.component';

@NgModule({
  declarations: [StringListComponent],
  imports: [
    CommonModule,
    InputModule,
    FormsModule,
    ReactiveFormsModule,
    TranslateModule,
    MatIconModule,
    MatLegacyTooltipModule,
    MatLegacyButtonModule,
  ],
  exports: [StringListComponent],
})
export class StringListModule {}
