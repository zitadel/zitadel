import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { MatButtonModule } from '@angular/material/button';
import { MatCheckboxModule } from '@angular/material/checkbox';
import { MatDialogModule } from '@angular/material/dialog';
import { MatIconModule } from '@angular/material/icon';
import { MatProgressSpinnerModule } from '@angular/material/progress-spinner';
import { MatTableModule } from '@angular/material/table';
import { MatTooltipModule } from '@angular/material/tooltip';
import { TranslateModule } from '@ngx-translate/core';
import { HasRoleModule } from 'src/app/directives/has-role/has-role.module';
import { AvatarModule } from 'src/app/modules/avatar/avatar.module';
import { CardModule } from 'src/app/modules/card/card.module';
import { InputModule } from 'src/app/modules/input/input.module';
import { PaginatorModule } from 'src/app/modules/paginator/paginator.module';
import { RefreshTableModule } from 'src/app/modules/refresh-table/refresh-table.module';
import { SharedModule } from 'src/app/modules/shared/shared.module';
import { HasRolePipeModule } from 'src/app/pipes/has-role-pipe/has-role-pipe.module';

import { UserListRoutingModule } from './user-list-routing.module';
import { UserListComponent } from './user-list.component';
import { UserTableComponent } from './user-table/user-table.component';


@NgModule({
    declarations: [
        UserListComponent,
        UserTableComponent,
    ],
    imports: [
        AvatarModule,
        UserListRoutingModule,
        CommonModule,
        FormsModule,
        MatButtonModule,
        MatDialogModule,
        HasRoleModule,
        CardModule,
        MatTableModule,
        MatIconModule,
        MatProgressSpinnerModule,
        MatCheckboxModule,
        MatTooltipModule,
        HasRolePipeModule,
        TranslateModule,
        SharedModule,
        RefreshTableModule,
        InputModule,
        PaginatorModule
    ],
    exports: [
        UserListComponent,
    ],
})
export class UserListModule { }
