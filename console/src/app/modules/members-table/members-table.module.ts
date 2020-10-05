import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { MatCheckboxModule } from '@angular/material/checkbox';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatIconModule } from '@angular/material/icon';
import { MatPaginatorModule } from '@angular/material/paginator';
import { MatSelectModule } from '@angular/material/select';
import { MatSortModule } from '@angular/material/sort';
import { MatTableModule } from '@angular/material/table';
import { MatTooltipModule } from '@angular/material/tooltip';
import { RouterModule } from '@angular/router';
import { TranslateModule } from '@ngx-translate/core';

import { AvatarModule } from '../avatar/avatar.module';
import { RefreshTableModule } from '../refresh-table/refresh-table.module';
import { MembersTableComponent } from './members-table.component';

@NgModule({
    declarations: [
        MembersTableComponent,
    ],
    imports: [
        CommonModule,
        MatFormFieldModule,
        MatSelectModule,
        MatCheckboxModule,
        MatIconModule,
        MatTableModule,
        MatPaginatorModule,
        MatSortModule,
        MatTooltipModule,
        FormsModule,
        TranslateModule,
        RefreshTableModule,
        RouterModule,
        AvatarModule,
    ],
    exports: [
        MembersTableComponent,
    ],
})
export class MembersTableModule { }
