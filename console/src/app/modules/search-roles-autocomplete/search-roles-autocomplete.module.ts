import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';
import { FormsModule, ReactiveFormsModule } from '@angular/forms';
import { MatIconModule } from '@angular/material/icon';
import { MatLegacyAutocompleteModule as MatAutocompleteModule } from '@angular/material/legacy-autocomplete';
import { MatLegacyButtonModule as MatButtonModule } from '@angular/material/legacy-button';
import { MatLegacyChipsModule as MatChipsModule } from '@angular/material/legacy-chips';
import { MatLegacyProgressSpinnerModule as MatProgressSpinnerModule } from '@angular/material/legacy-progress-spinner';
import { TranslateModule } from '@ngx-translate/core';
import { InputModule } from 'src/app/modules/input/input.module';

import { SearchRolesAutocompleteComponent } from './search-roles-autocomplete.component';

@NgModule({
  declarations: [SearchRolesAutocompleteComponent],
  imports: [
    CommonModule,
    MatAutocompleteModule,
    MatChipsModule,
    MatButtonModule,
    InputModule,
    MatIconModule,
    ReactiveFormsModule,
    MatProgressSpinnerModule,
    FormsModule,
    TranslateModule,
  ],
  exports: [SearchRolesAutocompleteComponent],
})
export class SearchRolesAutocompleteModule {}
