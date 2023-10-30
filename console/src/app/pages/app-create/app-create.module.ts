import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';
import { FormsModule, ReactiveFormsModule } from '@angular/forms';
import { MatButtonModule } from '@angular/material/button';
import { MatCheckboxModule } from '@angular/material/checkbox';
import { MatIconModule } from '@angular/material/icon';
import { MatSelectModule } from '@angular/material/select';
import { MatSlideToggleModule } from '@angular/material/slide-toggle';
import { TranslateModule } from '@ngx-translate/core';
import { HasRoleModule } from 'src/app/directives/has-role/has-role.module';
import { CreateLayoutModule } from 'src/app/modules/create-layout/create-layout.module';
import { InfoSectionModule } from 'src/app/modules/info-section/info-section.module';
import { InputModule } from 'src/app/modules/input/input.module';
import { PasswordComplexityViewModule } from 'src/app/modules/password-complexity-view/password-complexity-view.module';
import { HasRolePipeModule } from 'src/app/pipes/has-role-pipe/has-role-pipe.module';

import { SearchProjectAutocompleteModule } from 'src/app/modules/search-project-autocomplete/search-project-autocomplete.module';
import { AppCreateRoutingModule } from './app-create-routing.module';
import { AppCreateComponent } from './app-create.component';

@NgModule({
  declarations: [AppCreateComponent],
  imports: [
    AppCreateRoutingModule,
    CommonModule,
    FormsModule,
    ReactiveFormsModule,
    InfoSectionModule,
    InputModule,
    MatButtonModule,
    MatIconModule,
    SearchProjectAutocompleteModule,
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
export default class AppCreateModule {}
