import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';
import { MatButtonModule } from '@angular/material/button';
import { MatCheckboxModule } from '@angular/material/checkbox';
import { MatIconModule } from '@angular/material/icon';
import { MatMenuModule } from '@angular/material/menu';
import { MatPaginatorModule } from '@angular/material/paginator';
import { MatProgressSpinnerModule } from '@angular/material/progress-spinner';
import { MatTableModule } from '@angular/material/table';
import { MatTooltipModule } from '@angular/material/tooltip';
import { RouterModule } from '@angular/router';
import { TranslateModule } from '@ngx-translate/core';
import { HasRoleModule } from 'src/app/directives/has-role/has-role.module';
import { PipesModule } from 'src/app/pipes/pipes.module';

import { ProjectRolesComponent } from './project-roles.component';

@NgModule({
    declarations: [ProjectRolesComponent],
    imports: [
        CommonModule,
        MatButtonModule,
        HasRoleModule,
        MatTableModule,
        MatPaginatorModule,
        MatIconModule,
        MatProgressSpinnerModule,
        MatCheckboxModule,
        RouterModule,
        MatTooltipModule,
        PipesModule,
        TranslateModule,
        MatMenuModule,
    ],
    exports: [
        ProjectRolesComponent,
    ],
})
export class ProjectRolesModule { }
