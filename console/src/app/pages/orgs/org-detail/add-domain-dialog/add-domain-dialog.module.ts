import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { MatButtonModule } from '@angular/material/button';
import { TranslateModule } from '@ngx-translate/core';
import { FormFieldModule } from 'src/app/modules/form-field/form-field.module';

import { AddDomainDialogComponent } from './add-domain-dialog.component';

@NgModule({
    declarations: [AddDomainDialogComponent],
    imports: [
        CommonModule,
        TranslateModule,
        MatButtonModule,
        FormFieldModule,
        FormsModule,
    ],
})
export class AddDomainDialogModule { }
