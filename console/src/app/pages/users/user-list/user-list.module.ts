import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { MatButtonModule } from '@angular/material/button';
import { MatCheckboxModule } from '@angular/material/checkbox';
import { MatDialogModule } from '@angular/material/dialog';
import { MatIconModule } from '@angular/material/icon';
import { MatPaginatorModule } from '@angular/material/paginator';
import { MatProgressSpinnerModule } from '@angular/material/progress-spinner';
import { MatTableModule } from '@angular/material/table';
import { MatTooltipModule } from '@angular/material/tooltip';
import { TranslateModule } from '@ngx-translate/core';
import { HasRoleModule } from 'src/app/directives/has-role/has-role.module';
import { AvatarModule } from 'src/app/modules/avatar/avatar.module';
import { CardModule } from 'src/app/modules/card/card.module';

import { UserListRoutingModule } from './user-list-routing.module';
import { UserListComponent } from './user-list.component';



@NgModule({
    declarations: [
        UserListComponent,
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
        MatPaginatorModule,
        MatIconModule,
        MatProgressSpinnerModule,
        MatCheckboxModule,
        MatTooltipModule,
        TranslateModule,
    ],
    exports: [
        UserListComponent,
    ],
})
export class UserListModule { }
