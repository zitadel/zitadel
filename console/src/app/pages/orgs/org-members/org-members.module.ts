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
import { RouterModule } from '@angular/router';
import { TranslateModule } from '@ngx-translate/core';
import { HasRoleModule } from 'src/app/directives/has-role/has-role.module';

import { ProjectMemberCreateDialogModule } from '../../../modules/add-member-dialog/project-member-create-dialog.module';
import { OrgMembersComponent } from './org-members.component';


@NgModule({
    declarations: [OrgMembersComponent],
    imports: [
        CommonModule,
        FormsModule,
        ProjectMemberCreateDialogModule,
        HasRoleModule,
        MatButtonModule,
        MatDialogModule,
        MatTableModule,
        MatPaginatorModule,
        MatIconModule,
        RouterModule,
        MatProgressSpinnerModule,
        MatCheckboxModule,
        MatTooltipModule,
        TranslateModule,
    ],
    exports: [
        OrgMembersComponent,
    ],
})
export class OrgMembersModule { }
