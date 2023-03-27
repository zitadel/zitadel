import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';
import { FormsModule, ReactiveFormsModule } from '@angular/forms';
import { MatRippleModule } from '@angular/material/core';
import { MatIconModule } from '@angular/material/icon';
import { MatLegacyButtonModule as MatButtonModule } from '@angular/material/legacy-button';
import { MatLegacyCheckboxModule as MatCheckboxModule } from '@angular/material/legacy-checkbox';
import { MatLegacyProgressBarModule as MatProgressBarModule } from '@angular/material/legacy-progress-bar';
import { MatLegacyProgressSpinnerModule as MatProgressSpinnerModule } from '@angular/material/legacy-progress-spinner';
import { MatLegacySelectModule as MatSelectModule } from '@angular/material/legacy-select';
import { MatLegacyTooltipModule as MatTooltipModule } from '@angular/material/legacy-tooltip';
import { TranslateModule } from '@ngx-translate/core';
import { CreateLayoutModule } from 'src/app/modules/create-layout/create-layout.module';
import { DetailLayoutModule } from 'src/app/modules/detail-layout/detail-layout.module';
import { InfoSectionModule } from 'src/app/modules/info-section/info-section.module';
import { InputModule } from 'src/app/modules/input/input.module';
import { PasswordComplexityViewModule } from 'src/app/modules/password-complexity-view/password-complexity-view.module';

import { CountryCallingCodesService } from 'src/app/services/country-calling-codes.service';
import { UserCreateRoutingModule } from './user-create-routing.module';
import { UserCreateComponent } from './user-create.component';

@NgModule({
  declarations: [UserCreateComponent],
  providers: [CountryCallingCodesService],
  imports: [
    UserCreateRoutingModule,
    CommonModule,
    FormsModule,
    ReactiveFormsModule,
    MatSelectModule,
    CreateLayoutModule,
    MatButtonModule,
    MatIconModule,
    MatProgressSpinnerModule,
    MatProgressBarModule,
    PasswordComplexityViewModule,
    MatCheckboxModule,
    MatTooltipModule,
    TranslateModule,
    InfoSectionModule,
    DetailLayoutModule,
    InputModule,
    MatRippleModule,
  ],
})
export default class UserCreateModule {}
