import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';
import { FormsModule, ReactiveFormsModule } from '@angular/forms';
import { MatAutocompleteModule } from '@angular/material/autocomplete';
import { MatButtonModule } from '@angular/material/button';
import { MatCheckboxModule } from '@angular/material/checkbox';
import { MatChipsModule } from '@angular/material/chips';
import { MatDialogModule } from '@angular/material/dialog';
import { MatIconModule } from '@angular/material/icon';
import { MatPaginatorModule } from '@angular/material/paginator';
import { MatProgressSpinnerModule } from '@angular/material/progress-spinner';
import { MatSelectModule } from '@angular/material/select';
import { MatTableModule } from '@angular/material/table';
import { MatTooltipModule } from '@angular/material/tooltip';
import { TranslateModule } from '@ngx-translate/core';
import { HasRoleModule } from 'src/app/directives/has-role/has-role.module';
import { ActionKeysModule } from 'src/app/modules/action-keys/action-keys.module';
import { DetailLayoutModule } from 'src/app/modules/detail-layout/detail-layout.module';
import { InputModule } from 'src/app/modules/input/input.module';
import { MembersTableModule } from 'src/app/modules/members-table/members-table.module';
import { ProjectRoleChipModule } from 'src/app/modules/project-role-chip/project-role-chip.module';
import { HasRolePipeModule } from 'src/app/pipes/has-role-pipe/has-role-pipe.module';

import { ProjectGrantDetailRoutingModule } from './project-grant-detail-routing.module';
import { ProjectGrantDetailComponent } from './project-grant-detail.component';
import { ProjectGrantIllustrationComponent } from './project-grant-illustration/project-grant-illustration.component';
import {
  ProjectGrantMembersCreateDialogModule,
} from './project-grant-members-create-dialog/project-grant-members-create-dialog.module';

@NgModule({
  declarations: [ProjectGrantDetailComponent, ProjectGrantIllustrationComponent],
  imports: [
    CommonModule,
    ProjectGrantDetailRoutingModule,
    ProjectGrantMembersCreateDialogModule,
    MatAutocompleteModule,
    HasRoleModule,
    MatChipsModule,
    MatButtonModule,
    MatCheckboxModule,
    MatIconModule,
    MatTableModule,
    MatPaginatorModule,
    InputModule,
    MatTooltipModule,
    ReactiveFormsModule,
    MatProgressSpinnerModule,
    ActionKeysModule,
    ProjectRoleChipModule,
    FormsModule,
    TranslateModule,
    MatSelectModule,
    DetailLayoutModule,
    HasRolePipeModule,
    MembersTableModule,
    MatDialogModule,
  ],
})
export class ProjectGrantDetailModule {}
