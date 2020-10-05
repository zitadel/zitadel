import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';
import { ReactiveFormsModule } from '@angular/forms';
import { MatButtonModule } from '@angular/material/button';
import { MatChipsModule } from '@angular/material/chips';
import { MatIconModule } from '@angular/material/icon';
import { MatProgressSpinnerModule } from '@angular/material/progress-spinner';
import { MatTooltipModule } from '@angular/material/tooltip';
import { TranslateModule } from '@ngx-translate/core';
import { HasRoleModule } from 'src/app/directives/has-role/has-role.module';
import { DetailLayoutModule } from 'src/app/modules/detail-layout/detail-layout.module';
import { MembersTableModule } from 'src/app/modules/members-table/members-table.module';
import { RefreshTableModule } from 'src/app/modules/refresh-table/refresh-table.module';

import { OrgMembersRoutingModule } from './org-members-routing.module';
import { OrgMembersComponent } from './org-members.component';


@NgModule({
    declarations: [OrgMembersComponent],
    imports: [
        OrgMembersRoutingModule,
        CommonModule,
        MatChipsModule,
        MatButtonModule,
        HasRoleModule,
        MatIconModule,
        MatTooltipModule,
        ReactiveFormsModule,
        MatProgressSpinnerModule,
        TranslateModule,
        DetailLayoutModule,
        RefreshTableModule,
        MembersTableModule,
    ],
})
export class OrgMembersModule { }
