import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';
import { MatDialogModule } from '@angular/material/dialog';
import { RouterModule } from '@angular/router';
import { TranslateModule } from '@ngx-translate/core';
import { HasRoleModule } from 'src/app/directives/has-role/has-role.module';
import { ContributorsModule } from 'src/app/modules/contributors/contributors.module';

import { MemberCreateDialogModule } from '../../../modules/add-member-dialog/member-create-dialog.module';
import { OrgContributorsComponent } from './org-contributors.component';


@NgModule({
    declarations: [OrgContributorsComponent],
    imports: [
        CommonModule,
        MemberCreateDialogModule,
        HasRoleModule,
        MatDialogModule,
        RouterModule,
        TranslateModule,
        ContributorsModule,
    ],
    exports: [
        OrgContributorsComponent,
    ],
})
export class OrgContributorsModule { }
