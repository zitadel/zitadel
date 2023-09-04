import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';
import { FormsModule, ReactiveFormsModule } from '@angular/forms';
import { MatIconModule } from '@angular/material/icon';
import { MatButtonModule } from '@angular/material/button';
import { MatLegacyChipsModule } from '@angular/material/legacy-chips';
import { MatTooltipModule } from '@angular/material/tooltip';
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
    MatLegacyChipsModule,
    TranslateModule,
    MatIconModule,
    MatTooltipModule,
    MatButtonModule,
  ],
  exports: [StringListComponent],
})
export class StringListModule {}
