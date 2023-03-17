import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';
import { FormsModule, ReactiveFormsModule } from '@angular/forms';
import { MatIconModule } from '@angular/material/icon';
import { MatLegacyButtonModule as MatButtonModule } from '@angular/material/legacy-button';
import { MatLegacyCheckboxModule as MatCheckboxModule } from '@angular/material/legacy-checkbox';
import { MatLegacySelectModule as MatSelectModule } from '@angular/material/legacy-select';
import { MatLegacySlideToggleModule as MatSlideToggleModule } from '@angular/material/legacy-slide-toggle';
import { TranslateModule } from '@ngx-translate/core';
import { HasRoleModule } from 'src/app/directives/has-role/has-role.module';
import { CreateLayoutModule } from 'src/app/modules/create-layout/create-layout.module';
import { InfoSectionModule } from 'src/app/modules/info-section/info-section.module';
import { InputModule } from 'src/app/modules/input/input.module';
import { PasswordComplexityViewModule } from 'src/app/modules/password-complexity-view/password-complexity-view.module';
import { HasRolePipeModule } from 'src/app/pipes/has-role-pipe/has-role-pipe.module';

import { OrgCreateRoutingModule } from './org-create-routing.module';
import { OrgCreateComponent } from './org-create.component';

@NgModule({
  declarations: [OrgCreateComponent],
  imports: [
    OrgCreateRoutingModule,
    CommonModule,
    FormsModule,
    ReactiveFormsModule,
    InfoSectionModule,
    InputModule,
    MatButtonModule,
    MatIconModule,
    MatSelectModule,
    CreateLayoutModule,
    HasRolePipeModule,
    TranslateModule,
    HasRoleModule,
    MatCheckboxModule,
    PasswordComplexityViewModule,
    MatSlideToggleModule,
  ],
})
export default class OrgCreateModule {}
