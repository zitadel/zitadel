import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';
import { FormsModule, ReactiveFormsModule } from '@angular/forms';
import { MatButtonModule } from '@angular/material/button';
import { MatIconModule } from '@angular/material/icon';
import { MatTooltipModule } from '@angular/material/tooltip';
import { TranslateModule } from '@ngx-translate/core';
import { InputModule } from 'src/app/modules/input/input.module';

import { ProjectRoleCreateRoutingModule } from './project-role-create-routing.module';
import { ProjectRoleCreateComponent } from './project-role-create.component';

@NgModule({
    declarations: [ProjectRoleCreateComponent],
    imports: [
        ProjectRoleCreateRoutingModule,
        CommonModule,
        FormsModule,
        ReactiveFormsModule,
        InputModule,
        MatButtonModule,
        MatIconModule,
        MatTooltipModule,
        TranslateModule,
    ],
})
export class ProjectRoleCreateModule { }
