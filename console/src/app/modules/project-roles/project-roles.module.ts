import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';
import { FormsModule, ReactiveFormsModule } from '@angular/forms';
import { MatButtonModule } from '@angular/material/button';
import { MatCheckboxModule } from '@angular/material/checkbox';
import { MatDialogModule } from '@angular/material/dialog';
import { MatIconModule } from '@angular/material/icon';
import { MatMenuModule } from '@angular/material/menu';
import { MatProgressSpinnerModule } from '@angular/material/progress-spinner';
import { MatTableModule } from '@angular/material/table';
import { MatTooltipModule } from '@angular/material/tooltip';
import { RouterModule } from '@angular/router';
import { TranslateModule } from '@ngx-translate/core';
import { HasRoleModule } from 'src/app/directives/has-role/has-role.module';
import { InputModule } from 'src/app/modules/input/input.module';
import { HasRolePipeModule } from 'src/app/pipes/has-role-pipe/has-role-pipe.module';
import { LocalizedDatePipeModule } from 'src/app/pipes/localized-date-pipe/localized-date-pipe.module';
import { TimestampToDatePipeModule } from 'src/app/pipes/timestamp-to-date-pipe/timestamp-to-date-pipe.module';

import { PaginatorModule } from '../paginator/paginator.module';
import { RefreshTableModule } from '../refresh-table/refresh-table.module';
import { ProjectRoleDetailComponent } from './project-role-detail/project-role-detail.component';
import { ProjectRolesComponent } from './project-roles.component';


@NgModule({
    declarations: [ProjectRolesComponent, ProjectRoleDetailComponent],
    imports: [
        CommonModule,
        MatButtonModule,
        HasRoleModule,
        MatTableModule,
        PaginatorModule,
        MatDialogModule,
        InputModule,
        FormsModule,
        ReactiveFormsModule,
        MatIconModule,
        MatProgressSpinnerModule,
        MatCheckboxModule,
        RouterModule,
        MatTooltipModule,
        HasRolePipeModule,
        TranslateModule,
        MatMenuModule,
        TimestampToDatePipeModule,
        RefreshTableModule,
        LocalizedDatePipeModule,
    ],
    exports: [
        ProjectRolesComponent,
    ],
})
export class ProjectRolesModule { }
