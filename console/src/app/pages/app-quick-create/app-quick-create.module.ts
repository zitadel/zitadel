import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';
import { FormsModule, ReactiveFormsModule } from '@angular/forms';
import { MatButtonModule } from '@angular/material/button';
import { MatCheckboxModule } from '@angular/material/checkbox';
import { MatDialogModule } from '@angular/material/dialog';
import { MatIconModule } from '@angular/material/icon';
import { MatSelectModule } from '@angular/material/select';
import { MatTooltipModule } from '@angular/material/tooltip';
import { TranslateModule } from '@ngx-translate/core';
import { HasRoleModule } from 'src/app/directives/has-role/has-role.module';
import { CopyToClipboardModule } from 'src/app/directives/copy-to-clipboard/copy-to-clipboard.module';
import { CreateLayoutModule } from 'src/app/modules/create-layout/create-layout.module';
import { InfoSectionModule } from 'src/app/modules/info-section/info-section.module';
import { InputModule } from 'src/app/modules/input/input.module';
import { PasswordComplexityViewModule } from 'src/app/modules/password-complexity-view/password-complexity-view.module';
import { HasRolePipeModule } from 'src/app/pipes/has-role-pipe/has-role-pipe.module';

import { SearchProjectAutocompleteModule } from 'src/app/modules/search-project-autocomplete/search-project-autocomplete.module';
import { AppCreateRoutingModule } from './app-quick-create-routing.module';
import { AppQuickCreateComponent } from './app-quick-create.component';
import { FrameworkAutocompleteComponent } from 'src/app/components/framework-autocomplete/framework-autocomplete.component';
import { FrameworkChangeComponent } from 'src/app/components/framework-change/framework-change.component';

@NgModule({
  declarations: [AppQuickCreateComponent],
  imports: [
    FrameworkChangeComponent,
    AppCreateRoutingModule,
    FrameworkAutocompleteComponent,
    CommonModule,
    FormsModule,
    ReactiveFormsModule,
    InfoSectionModule,
    InputModule,
    MatButtonModule,
    MatDialogModule,
    MatIconModule,
    MatTooltipModule,
    SearchProjectAutocompleteModule,
    MatSelectModule,
    CreateLayoutModule,
    HasRolePipeModule,
    TranslateModule,
    HasRoleModule,
    CopyToClipboardModule,
  ],
})
export default class AppCreateModule {}
