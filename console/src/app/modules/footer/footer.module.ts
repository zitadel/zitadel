import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';
import { TranslateModule } from '@ngx-translate/core';

import { ThemeSettingModule } from '../theme-setting/theme-setting.module';
import { FooterComponent } from './footer.component';

@NgModule({
  declarations: [FooterComponent],
  imports: [CommonModule, TranslateModule, ThemeSettingModule],
  exports: [FooterComponent],
})
export class FooterModule {}
