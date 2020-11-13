import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';
import { FormsModule, ReactiveFormsModule } from '@angular/forms';
import { MatAutocompleteModule } from '@angular/material/autocomplete';
import { MatButtonModule } from '@angular/material/button';
import { MatChipsModule } from '@angular/material/chips';
import { MatIconModule } from '@angular/material/icon';
import { MatProgressSpinnerModule } from '@angular/material/progress-spinner';
import { TranslateModule } from '@ngx-translate/core';
import { FormFieldModule } from 'src/app/modules/form-field/form-field.module';

import { SearchRolesAutocompleteComponent } from './search-roles-autocomplete.component';


@NgModule({
    declarations: [SearchRolesAutocompleteComponent],
    imports: [
        CommonModule,
        MatAutocompleteModule,
        MatChipsModule,
        MatButtonModule,
        FormFieldModule,
        MatIconModule,
        ReactiveFormsModule,
        MatProgressSpinnerModule,
        FormsModule,
        TranslateModule,
    ],
    exports: [
        SearchRolesAutocompleteComponent,
    ],
})
export class SearchRolesAutocompleteModule { }
