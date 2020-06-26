import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';
import { MatButtonModule } from '@angular/material/button';
import { MatIconModule } from '@angular/material/icon';
import { MatTooltipModule } from '@angular/material/tooltip';
import { TranslateModule } from '@ngx-translate/core';

import { MemberCreateDialogModule } from '../add-member-dialog/member-create-dialog.module';
import { AvatarModule } from '../avatar/avatar.module';
import { ProjectContributorsComponent } from './project-contributors.component';



@NgModule({
    declarations: [ProjectContributorsComponent],
    imports: [
        MemberCreateDialogModule,
        CommonModule,
        TranslateModule,
        MatTooltipModule,
        MatIconModule,
        MatButtonModule,
        AvatarModule,
    ],
    exports: [
        ProjectContributorsComponent,
    ],
})
export class ProjectContributorsModule { }

