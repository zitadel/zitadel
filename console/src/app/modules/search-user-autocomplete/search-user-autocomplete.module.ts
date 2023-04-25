import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';
import { FormsModule, ReactiveFormsModule } from '@angular/forms';
import { MatIconModule } from '@angular/material/icon';
import { MatLegacyAutocompleteModule as MatAutocompleteModule } from '@angular/material/legacy-autocomplete';
import { MatLegacyButtonModule as MatButtonModule } from '@angular/material/legacy-button';
import { MatLegacyChipsModule as MatChipsModule } from '@angular/material/legacy-chips';
import { MatLegacyProgressSpinnerModule as MatProgressSpinnerModule } from '@angular/material/legacy-progress-spinner';
import { MatLegacySelectModule as MatSelectModule } from '@angular/material/legacy-select';
import { MatLegacyTooltipModule as MatTooltipModule } from '@angular/material/legacy-tooltip';
import { TranslateModule } from '@ngx-translate/core';

import { AvatarModule } from '../avatar/avatar.module';
import { InputModule } from '../input/input.module';
import { SearchUserAutocompleteComponent } from './search-user-autocomplete.component';

@NgModule({
  declarations: [SearchUserAutocompleteComponent],
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
    MatTooltipModule,
    TranslateModule,
    MatSelectModule,
    AvatarModule,
  ],
  exports: [SearchUserAutocompleteComponent],
})
export class SearchUserAutocompleteModule {}
