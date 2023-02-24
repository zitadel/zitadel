import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';
import { MatIconModule } from '@angular/material/icon';
import { MatLegacyButtonModule as MatButtonModule } from '@angular/material/legacy-button';
import { TranslateModule } from '@ngx-translate/core';

import { InputModule } from '../input/input.module';
import { ProjectRolesTableModule } from '../project-roles-table/project-roles-table.module';
import { UserGrantRoleDialogComponent } from './user-grant-role-dialog.component';

@NgModule({
  declarations: [UserGrantRoleDialogComponent],
  imports: [CommonModule, InputModule, MatButtonModule, MatIconModule, TranslateModule, ProjectRolesTableModule],
})
export class UserGrantRoleDialogModule {}
