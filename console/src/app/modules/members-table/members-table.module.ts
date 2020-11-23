import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { MatButtonModule } from '@angular/material/button';
import { MatCheckboxModule } from '@angular/material/checkbox';
import { MatIconModule } from '@angular/material/icon';
import { MatPaginatorModule } from '@angular/material/paginator';
import { MatSelectModule } from '@angular/material/select';
import { MatSortModule } from '@angular/material/sort';
import { MatTableModule } from '@angular/material/table';
import { MatTooltipModule } from '@angular/material/tooltip';
import { RouterModule } from '@angular/router';
import { TranslateModule } from '@ngx-translate/core';
import { InputModule } from 'src/app/modules/input/input.module';

import { AvatarModule } from '../avatar/avatar.module';
import { RefreshTableModule } from '../refresh-table/refresh-table.module';
import { MembersTableComponent } from './members-table.component';

@NgModule({
    declarations: [
        MembersTableComponent,
    ],
    imports: [
        CommonModule,
        InputModule,
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
        MatButtonModule,
    ],
    exports: [
        MembersTableComponent,
    ],
})
export class MembersTableModule { }
