import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { MatButtonModule } from '@angular/material/button';
import { MatDialogModule } from '@angular/material/dialog';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatIconModule } from '@angular/material/icon';
import { MatInputModule } from '@angular/material/input';
import { TranslateModule } from '@ngx-translate/core';

import { CodeDialogComponent } from './code-dialog.component';

@NgModule({
    declarations: [
        CodeDialogComponent,
    ],
    imports: [
        CommonModule,
        FormsModule,
        MatDialogModule,
        MatFormFieldModule,
        MatInputModule,
        MatButtonModule,
        MatIconModule,
        TranslateModule,
    ],
    entryComponents: [
        CodeDialogComponent,
    ],
    exports: [
        CodeDialogComponent,
    ],
})
export class CodeDialogModule { }
