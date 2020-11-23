import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { MatButtonModule } from '@angular/material/button';
import { MatDialogModule } from '@angular/material/dialog';
import { MatSelectModule } from '@angular/material/select';
import { TranslateModule } from '@ngx-translate/core';
import { InputModule } from 'src/app/modules/input/input.module';

import { AddIdpDialogComponent } from './add-idp-dialog.component';

@NgModule({
    declarations: [AddIdpDialogComponent],
    imports: [
        CommonModule,
        MatDialogModule,
        MatButtonModule,
        TranslateModule,
        InputModule,
        MatSelectModule,
        FormsModule,
    ],
})
export class AddIdpDialogModule { }
