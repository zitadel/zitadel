import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { MatButtonModule } from '@angular/material/button';
import { MatCheckboxModule } from '@angular/material/checkbox';
import { MatIconModule } from '@angular/material/icon';
import { MatPaginatorModule } from '@angular/material/paginator';
import { MatProgressSpinnerModule } from '@angular/material/progress-spinner';
import { MatTableModule } from '@angular/material/table';
import { MatTooltipModule } from '@angular/material/tooltip';
import { RouterModule } from '@angular/router';
import { TranslateModule } from '@ngx-translate/core';
import { HasRoleModule } from 'src/app/directives/has-role/has-role.module';
import { PipesModule } from 'src/app/pipes/pipes.module';

import { AvatarModule } from '../avatar/avatar.module';
import { UserGrantsComponent } from './user-grants.component';



@NgModule({
    declarations: [UserGrantsComponent],
    imports: [
        CommonModule,
        FormsModule,
        AvatarModule,
        MatButtonModule,
        HasRoleModule,
        MatTableModule,
        MatPaginatorModule,
        MatIconModule,
        RouterModule,
        MatProgressSpinnerModule,
        MatCheckboxModule,
        MatTooltipModule,
        TranslateModule,
        PipesModule,
    ],
    exports: [
        UserGrantsComponent,
    ],
})
export class UserGrantsModule { }
