import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';
import { FormsModule, ReactiveFormsModule } from '@angular/forms';
import { MatButtonModule } from '@angular/material/button';
import { MatCheckboxModule } from '@angular/material/checkbox';
import { MatRippleModule } from '@angular/material/core';
import { MatDialogModule } from '@angular/material/dialog';
import { MatIconModule } from '@angular/material/icon';
import { MatProgressSpinnerModule } from '@angular/material/progress-spinner';
import { MatSelectModule } from '@angular/material/select';
import { MatSlideToggleModule } from '@angular/material/slide-toggle';
import { MatTooltipModule } from '@angular/material/tooltip';
import { TranslateModule } from '@ngx-translate/core';
import { HasRoleModule } from 'src/app/directives/has-role/has-role.module';
import { CardModule } from 'src/app/modules/card/card.module';
import { DetailLayoutModule } from 'src/app/modules/detail-layout/detail-layout.module';
import { InputModule } from 'src/app/modules/input/input.module';
import { HasRolePipeModule } from 'src/app/pipes/has-role-pipe/has-role-pipe.module';

import { InfoSectionModule } from '../../info-section/info-section.module';
import { WarnDialogModule } from '../../warn-dialog/warn-dialog.module';
import { DialogAddTypeComponent } from './factor-table/dialog-add-type/dialog-add-type.component';
import { FactorTableComponent } from './factor-table/factor-table.component';
import { LoginPolicyRoutingModule } from './login-policy-routing.module';
import { LoginPolicyComponent } from './login-policy.component';

@NgModule({
  declarations: [LoginPolicyComponent, FactorTableComponent, DialogAddTypeComponent],
  imports: [
    LoginPolicyRoutingModule,
    CommonModule,
    InfoSectionModule,
    FormsModule,
    CardModule,
    ReactiveFormsModule,
    InputModule,
    MatIconModule,
    MatButtonModule,
    MatSlideToggleModule,
    WarnDialogModule,
    HasRoleModule,
    MatDialogModule,
    HasRolePipeModule,
    MatCheckboxModule,
    MatTooltipModule,
    DetailLayoutModule,
    MatProgressSpinnerModule,
    MatSelectModule,
    MatRippleModule,
    TranslateModule,
  ],
  exports: [LoginPolicyComponent],
})
export class LoginPolicyModule {}
