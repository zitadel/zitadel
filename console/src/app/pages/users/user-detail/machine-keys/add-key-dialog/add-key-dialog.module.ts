import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { MatButtonModule } from '@angular/material/button';
import { MatDatepickerModule } from '@angular/material/datepicker';
import { MatIconModule } from '@angular/material/icon';
import { MatSelectModule } from '@angular/material/select';
import { TranslateModule } from '@ngx-translate/core';
import { FormFieldModule } from 'src/app/modules/form-field/form-field.module';

import { AddKeyDialogComponent } from './add-key-dialog.component';

@NgModule({
    declarations: [AddKeyDialogComponent],
    imports: [
        CommonModule,
        TranslateModule,
        MatButtonModule,
        FormFieldModule,
        MatSelectModule,
        MatIconModule,
        FormsModule,
        MatDatepickerModule,
    ],
})
export class AddKeyDialogModule { }
