import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';
import { MatButtonModule } from '@angular/material/button';
import { MatDialogModule } from '@angular/material/dialog';
import { MatIconModule } from '@angular/material/icon';
import { MatTooltipModule } from '@angular/material/tooltip';
import { TranslateModule } from '@ngx-translate/core';
import { HasRoleModule } from 'src/app/directives/has-role/has-role.module';
import { DetailLayoutModule } from 'src/app/modules/detail-layout/detail-layout.module';
import { MembersTableModule } from 'src/app/modules/members-table/members-table.module';

import { ProjectMembersRoutingModule } from './project-members-routing.module';
import { ProjectMembersComponent } from './project-members.component';

@NgModule({
    declarations: [ProjectMembersComponent],
    imports: [
        ProjectMembersRoutingModule,
        CommonModule,
        HasRoleModule,
        MatButtonModule,
        MatIconModule,
        MatTooltipModule,
        TranslateModule,
        DetailLayoutModule,
        MatDialogModule,
        MembersTableModule,
    ],
})
export class ProjectMembersModule { }
