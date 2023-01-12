import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { MatLegacyButtonModule as MatButtonModule } from '@angular/material/legacy-button';
import { MatLegacyMenuModule as MatMenuModule } from '@angular/material/legacy-menu';
import { TranslateModule } from '@ngx-translate/core';

import { ThemeSettingComponent } from './theme-setting.component';

@NgModule({
  declarations: [ThemeSettingComponent],
  imports: [CommonModule, FormsModule, MatButtonModule, MatMenuModule, TranslateModule],
  exports: [ThemeSettingComponent],
})
export class ThemeSettingModule {}
