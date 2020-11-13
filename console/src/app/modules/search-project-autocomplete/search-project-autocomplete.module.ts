import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';
import { FormsModule, ReactiveFormsModule } from '@angular/forms';
import { MatAutocompleteModule } from '@angular/material/autocomplete';
import { MatButtonModule } from '@angular/material/button';
import { MatChipsModule } from '@angular/material/chips';
import { MatIconModule } from '@angular/material/icon';
import { MatProgressSpinnerModule } from '@angular/material/progress-spinner';
import { MatSelectModule } from '@angular/material/select';
import { TranslateModule } from '@ngx-translate/core';
import { FormFieldModule } from 'src/app/modules/form-field/form-field.module';

import { SearchProjectAutocompleteComponent } from './search-project-autocomplete.component';

@NgModule({
    declarations: [
        SearchProjectAutocompleteComponent,
    ],
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
        MatSelectModule,
    ],
    exports: [
        SearchProjectAutocompleteComponent,
    ],
})
export class SearchProjectAutocompleteModule { }
