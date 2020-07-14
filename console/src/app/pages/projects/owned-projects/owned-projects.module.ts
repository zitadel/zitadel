import { CommonModule } from '@angular/common';
import { NgModule, NO_ERRORS_SCHEMA } from '@angular/core';
import { FormsModule, ReactiveFormsModule } from '@angular/forms';
import { MatButtonModule } from '@angular/material/button';
import { MatCheckboxModule } from '@angular/material/checkbox';
import { MatChipsModule } from '@angular/material/chips';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatIconModule } from '@angular/material/icon';
import { MatInputModule } from '@angular/material/input';
import { MatPaginatorModule } from '@angular/material/paginator';
import { MatProgressBarModule } from '@angular/material/progress-bar';
import { MatProgressSpinnerModule } from '@angular/material/progress-spinner';
import { MatSortModule } from '@angular/material/sort';
import { MatTableModule } from '@angular/material/table';
import { MatTooltipModule } from '@angular/material/tooltip';
import { TranslateModule } from '@ngx-translate/core';
import { HasRoleModule } from 'src/app/directives/has-role/has-role.module';
import { AvatarModule } from 'src/app/modules/avatar/avatar.module';
import { CardModule } from 'src/app/modules/card/card.module';
import { ChangesModule } from 'src/app/modules/changes/changes.module';
import { UserGrantsModule } from 'src/app/modules/user-grants/user-grants.module';
import { HasRolePipeModule } from 'src/app/pipes/has-role-pipe.module';
import { TimestampToDatePipeModule } from 'src/app/pipes/timestamp-to-date-pipe.module';

import { OwnedProjectGridComponent } from './owned-project-list/owned-project-grid/owned-project-grid.component';
import { OwnedProjectListComponent } from './owned-project-list/owned-project-list.component';
import { OwnedProjectsRoutingModule } from './owned-projects-routing.module';
import { OwnedProjectsComponent } from './owned-projects.component';
import { ProjectGrantsComponent } from './project-grants/project-grants.component';

@NgModule({
    declarations: [
        OwnedProjectsComponent,
        OwnedProjectListComponent,
        OwnedProjectGridComponent,
        ProjectGrantsComponent,
    ],
    imports: [
        CommonModule,
        OwnedProjectsRoutingModule,
        UserGrantsModule,
        FormsModule,
        ReactiveFormsModule,
        TranslateModule,
        AvatarModule,
        ReactiveFormsModule,
        HasRoleModule,
        MatTableModule,
        MatPaginatorModule,
        MatFormFieldModule,
        MatInputModule,
        ChangesModule,
        MatChipsModule,
        MatIconModule,
        MatButtonModule,
        MatProgressSpinnerModule,
        MatProgressBarModule,
        MatCheckboxModule,
        CardModule,
        MatTooltipModule,
        MatSortModule,
        HasRolePipeModule,
        TimestampToDatePipeModule,
    ],
    schemas: [NO_ERRORS_SCHEMA],
})
export class OwnedProjectsModule { }
