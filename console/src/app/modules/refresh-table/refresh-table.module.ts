import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { MatButtonModule } from '@angular/material/button';
import { MatIconModule } from '@angular/material/icon';
import { TranslateModule } from '@ngx-translate/core';

import { RefreshTableComponent } from './refresh-table.component';



@NgModule({
    declarations: [RefreshTableComponent],
    imports: [
        CommonModule,
        MatButtonModule,
        MatIconModule,
        TranslateModule,
        FormsModule,
    ],
    exports: [
        RefreshTableComponent,
    ],
})
export class RefreshTableModule { }
