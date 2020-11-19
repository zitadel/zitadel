import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { MatButtonModule } from '@angular/material/button';
import { TranslateModule } from '@ngx-translate/core';
import { InputModule } from 'src/app/modules/input/input.module';

import { AddDomainDialogComponent } from './add-domain-dialog.component';

@NgModule({
    declarations: [AddDomainDialogComponent],
    imports: [
        CommonModule,
        TranslateModule,
        MatButtonModule,
        InputModule,
        FormsModule,
    ],
})
export class AddDomainDialogModule { }
