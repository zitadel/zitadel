import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';
import { TranslateModule } from '@ngx-translate/core';

import { ThemeSettingModule } from '../theme-setting/theme-setting.module';
import { FooterComponent } from './footer.component';
import { FontAwesomeModule } from '@fortawesome/angular-fontawesome';

@NgModule({
  declarations: [FooterComponent],
  imports: [CommonModule, TranslateModule, ThemeSettingModule, FontAwesomeModule],
  exports: [FooterComponent],
})
export class FooterModule {}
