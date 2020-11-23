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
        TranslateModule,
        MatSelectModule,
        AvatarModule,
    ],
    exports: [SearchUserAutocompleteComponent],
})
export class SearchUserAutocompleteModule { }
