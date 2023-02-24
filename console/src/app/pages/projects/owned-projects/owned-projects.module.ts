import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';
import { FormsModule, ReactiveFormsModule } from '@angular/forms';
import { MatRippleModule } from '@angular/material/core';
import { MatIconModule } from '@angular/material/icon';
import { MatLegacyButtonModule as MatButtonModule } from '@angular/material/legacy-button';
import { MatLegacyCheckboxModule as MatCheckboxModule } from '@angular/material/legacy-checkbox';
import { MatLegacyChipsModule as MatChipsModule } from '@angular/material/legacy-chips';
import { MatLegacyProgressBarModule as MatProgressBarModule } from '@angular/material/legacy-progress-bar';
import { MatLegacyProgressSpinnerModule as MatProgressSpinnerModule } from '@angular/material/legacy-progress-spinner';
import { MatLegacyTableModule as MatTableModule } from '@angular/material/legacy-table';
import { MatLegacyTooltipModule as MatTooltipModule } from '@angular/material/legacy-tooltip';
import { MatSortModule } from '@angular/material/sort';
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
