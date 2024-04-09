import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';
import { FormsModule, ReactiveFormsModule } from '@angular/forms';
import { MatCheckboxModule } from '@angular/material/checkbox';
import { MatSelectModule } from '@angular/material/select';
import { TranslateModule } from '@ngx-translate/core';
import { InfoSectionModule } from '../info-section/info-section.module';
import { ProviderOptionsComponent } from './provider-options.component';

@NgModule({
  declarations: [ProviderOptionsComponent],
  imports: [
    CommonModule,
    MatCheckboxModule,
    MatSelectModule,
    FormsModule,
    InfoSectionModule,
    ReactiveFormsModule,
    TranslateModule,
  ],
  exports: [ProviderOptionsComponent],
})
export class ProviderOptionsModule {}
