import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { MatButtonModule } from '@angular/material/button';
import { MatDialogModule } from '@angular/material/dialog';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatSelectModule } from '@angular/material/select';
import { TranslateModule } from '@ngx-translate/core';

import { AddIdpDialogComponent } from './add-idp-dialog.component';

@NgModule({
    declarations: [AddIdpDialogComponent],
    imports: [
        CommonModule,
        MatDialogModule,
        MatButtonModule,
        TranslateModule,
        MatFormFieldModule,
        MatSelectModule,
        FormsModule,
    ],
})
export class AddIdpDialogModule { }
