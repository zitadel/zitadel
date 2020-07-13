import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';
import { MatButtonModule } from '@angular/material/button';
import { TranslateModule } from '@ngx-translate/core';

import { WarnDialogComponent } from './warn-dialog.component';



@NgModule({
    declarations: [WarnDialogComponent],
    imports: [
        CommonModule,
        TranslateModule,
        MatButtonModule,
    ],
})
export class WarnDialogModule { }
