import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';
import { MatButtonModule } from '@angular/material/button';
import { MatIconModule } from '@angular/material/icon';
import { TranslateModule } from '@ngx-translate/core';

import { MatDialogModule } from '@angular/material/dialog';
import { InputModule } from '../input/input.module';
import { ProjectRolesTableModule } from '../project-roles-table/project-roles-table.module';
import { UserGrantRoleDialogComponent } from './user-grant-role-dialog.component';

@NgModule({
  declarations: [UserGrantRoleDialogComponent],
  imports: [
    CommonModule,
    InputModule,
    MatDialogModule,
    MatButtonModule,
    MatIconModule,
    TranslateModule,
    ProjectRolesTableModule,
  ],
})
export class UserGrantRoleDialogModule {}
