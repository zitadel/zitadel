import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';
import { FormsModule, ReactiveFormsModule } from '@angular/forms';
import { MatLegacyCheckboxModule as MatCheckboxModule } from '@angular/material/legacy-checkbox';
import { TranslateModule } from '@ngx-translate/core';
import { InfoSectionModule } from '../info-section/info-section.module';
import { ProviderOptionsComponent } from './provider-options.component';

@NgModule({
  declarations: [ProviderOptionsComponent],
  imports: [CommonModule, MatCheckboxModule, FormsModule, InfoSectionModule, ReactiveFormsModule, TranslateModule],
  exports: [ProviderOptionsComponent],
})
export class ProviderOptionsModule {}
