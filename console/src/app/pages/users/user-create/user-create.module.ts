import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';
import { FormsModule, ReactiveFormsModule } from '@angular/forms';
import { MatButtonModule } from '@angular/material/button';
import { MatCheckboxModule } from '@angular/material/checkbox';
import { MatRippleModule } from '@angular/material/core';
import { MatIconModule } from '@angular/material/icon';
import { MatProgressBarModule } from '@angular/material/progress-bar';
import { MatProgressSpinnerModule } from '@angular/material/progress-spinner';
import { MatSelectModule } from '@angular/material/select';
import { MatTooltipModule } from '@angular/material/tooltip';
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
