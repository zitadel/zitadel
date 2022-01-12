import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { MatButtonModule } from '@angular/material/button';
import { MatMenuModule } from '@angular/material/menu';
import { TranslateModule } from '@ngx-translate/core';

import { ThemeSettingComponent } from './theme-setting.component';

@NgModule({
  declarations: [ThemeSettingComponent],
  imports: [CommonModule, FormsModule, MatButtonModule, MatMenuModule, TranslateModule],
  exports: [ThemeSettingComponent],
})
export class ThemeSettingModule {}
