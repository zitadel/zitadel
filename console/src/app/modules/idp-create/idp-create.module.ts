import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';
import { FormsModule, ReactiveFormsModule } from '@angular/forms';
import { MatButtonModule } from '@angular/material/button';
import { MatChipsModule } from '@angular/material/chips';
import { MatIconModule } from '@angular/material/icon';
import { MatProgressBarModule } from '@angular/material/progress-bar';
import { MatSelectModule } from '@angular/material/select';
import { MatTooltipModule } from '@angular/material/tooltip';
import { TranslateModule } from '@ngx-translate/core';
import { FormFieldModule } from 'src/app/modules/form-field/form-field.module';

import { IdpCreateRoutingModule } from './idp-create-routing.module';
import { IdpCreateComponent } from './idp-create.component';

@NgModule({
    declarations: [IdpCreateComponent],
    imports: [
        IdpCreateRoutingModule,
        CommonModule,
        FormsModule,
        ReactiveFormsModule,
        FormFieldModule,
        MatButtonModule,
        MatSelectModule,
        MatIconModule,
        MatChipsModule,
        MatTooltipModule,
        TranslateModule,
        MatProgressBarModule,
    ],
})
export class IdpCreateModule { }
