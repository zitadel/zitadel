import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';
import { FormsModule, ReactiveFormsModule } from '@angular/forms';
import { MatButtonModule } from '@angular/material/button';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatIconModule } from '@angular/material/icon';
import { MatInputModule } from '@angular/material/input';
import { MatTooltipModule } from '@angular/material/tooltip';
import { TranslateModule } from '@ngx-translate/core';

import { ProjectRoleCreateRoutingModule } from './project-role-create-routing.module';
import { ProjectRoleCreateComponent } from './project-role-create.component';

@NgModule({
    declarations: [ProjectRoleCreateComponent],
    imports: [
        ProjectRoleCreateRoutingModule,
        CommonModule,
        FormsModule,
        ReactiveFormsModule,
        MatInputModule,
        MatFormFieldModule,
        MatButtonModule,
        MatIconModule,
        MatTooltipModule,
        TranslateModule,
    ],
})
export class ProjectRoleCreateModule { }
