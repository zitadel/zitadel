import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { MatButtonModule } from '@angular/material/button';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatInputModule } from '@angular/material/input';
import { TranslateModule } from '@ngx-translate/core';

import { AddDomainDialogComponent } from './add-domain-dialog.component';

@NgModule({
    declarations: [AddDomainDialogComponent],
    imports: [
        CommonModule,
        TranslateModule,
        MatButtonModule,
        MatFormFieldModule,
        MatInputModule,
        FormsModule,
    ],
})
export class AddDomainDialogModule { }
