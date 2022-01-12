import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { MatMenuModule } from '@angular/material/menu';

import { ThemeSettingComponent } from './theme-setting.component';

@NgModule({
  declarations: [ThemeSettingComponent],
  imports: [CommonModule, FormsModule, MatMenuModule],
  exports: [ThemeSettingComponent],
})
export class ThemeSettingModule {}
