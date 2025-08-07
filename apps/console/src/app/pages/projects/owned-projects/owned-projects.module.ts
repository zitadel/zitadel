import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';
import { FormsModule, ReactiveFormsModule } from '@angular/forms';
import { MatButtonModule } from '@angular/material/button';
import { MatCheckboxModule } from '@angular/material/checkbox';
import { MatChipsModule } from '@angular/material/chips';
import { MatRippleModule } from '@angular/material/core';
import { MatIconModule } from '@angular/material/icon';
import { MatProgressBarModule } from '@angular/material/progress-bar';
import { MatProgressSpinnerModule } from '@angular/material/progress-spinner';
import { MatSortModule } from '@angular/material/sort';
import { MatTableModule } from '@angular/material/table';
import { MatTooltipModule } from '@angular/material/tooltip';
import { TranslateModule } from '@ngx-translate/core';
import { HasRoleModule } from 'src/app/directives/has-role/has-role.module';
import { AvatarModule } from 'src/app/modules/avatar/avatar.module';
import { CardModule } from 'src/app/modules/card/card.module';
import { InputModule } from 'src/app/modules/input/input.module';
import { PaginatorModule } from 'src/app/modules/paginator/paginator.module';
import { RefreshTableModule } from 'src/app/modules/refresh-table/refresh-table.module';
import { UserGrantsModule } from 'src/app/modules/user-grants/user-grants.module';
import { WarnDialogModule } from 'src/app/modules/warn-dialog/warn-dialog.module';
import { HasRolePipeModule } from 'src/app/pipes/has-role-pipe/has-role-pipe.module';
import { LocalizedDatePipeModule } from 'src/app/pipes/localized-date-pipe/localized-date-pipe.module';
import { TimestampToDatePipeModule } from 'src/app/pipes/timestamp-to-date-pipe/timestamp-to-date-pipe.module';

import { NameDialogModule } from '../../../modules/name-dialog/name-dialog.module';
import { ProjectPrivateLabelingDialogModule } from '../../../modules/project-private-labeling-dialog/project-private-labeling-dialog.module';
import { OwnedProjectsRoutingModule } from './owned-projects-routing.module';

@NgModule({
  declarations: [],
  imports: [
    CommonModule,
    OwnedProjectsRoutingModule,
    UserGrantsModule,
    FormsModule,
    NameDialogModule,
    TranslateModule,
    AvatarModule,
    ReactiveFormsModule,
    HasRoleModule,
    MatTableModule,
    PaginatorModule,
    ProjectPrivateLabelingDialogModule,
    InputModule,
    MatChipsModule,
    MatIconModule,
    WarnDialogModule,
    MatButtonModule,
    MatProgressSpinnerModule,
    MatProgressBarModule,
    MatCheckboxModule,
    CardModule,
    MatTooltipModule,
    MatSortModule,
    HasRolePipeModule,
    TimestampToDatePipeModule,
    LocalizedDatePipeModule,
    RefreshTableModule,
    MatRippleModule,
  ],
})
export default class OwnedProjectsModule {}
