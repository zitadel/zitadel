import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';
import { TranslateModule } from '@ngx-translate/core';
import { ContributorsModule } from 'src/app/modules/contributors/contributors.module';

import { MemberCreateDialogModule } from '../add-member-dialog/member-create-dialog.module';
import { ProjectContributorsComponent } from './project-contributors.component';


@NgModule({
    declarations: [ProjectContributorsComponent],
    imports: [
        MemberCreateDialogModule,
        CommonModule,
        TranslateModule,
        ContributorsModule,
    ],
    exports: [
        ProjectContributorsComponent,
    ],
})
export class ProjectContributorsModule { }

