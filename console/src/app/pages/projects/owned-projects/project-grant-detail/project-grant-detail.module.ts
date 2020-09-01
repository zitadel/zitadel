import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';
import { FormsModule, ReactiveFormsModule } from '@angular/forms';
import { MatAutocompleteModule } from '@angular/material/autocomplete';
import { MatButtonModule } from '@angular/material/button';
import { MatCheckboxModule } from '@angular/material/checkbox';
import { MatChipsModule } from '@angular/material/chips';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatIconModule } from '@angular/material/icon';
import { MatPaginatorModule } from '@angular/material/paginator';
import { MatProgressSpinnerModule } from '@angular/material/progress-spinner';
import { MatSelectModule } from '@angular/material/select';
import { MatTableModule } from '@angular/material/table';
import { MatTooltipModule } from '@angular/material/tooltip';
import { TranslateModule } from '@ngx-translate/core';
import { HasRoleModule } from 'src/app/directives/has-role/has-role.module';
import { DetailLayoutModule } from 'src/app/modules/detail-layout/detail-layout.module';
import { HasRolePipeModule } from 'src/app/pipes/has-role-pipe.module';

import { ProjectGrantDetailRoutingModule } from './project-grant-detail-routing.module';
import { ProjectGrantDetailComponent } from './project-grant-detail.component';
import { ProjectGrantMembersModule } from './project-grant-members/project-grant-members.module';


@NgModule({
    declarations: [ProjectGrantDetailComponent],
    imports: [
        CommonModule,
        ProjectGrantDetailRoutingModule,
        ProjectGrantMembersModule,
        MatAutocompleteModule,
        HasRoleModule,
        MatChipsModule,
        MatButtonModule,
        MatCheckboxModule,
        MatIconModule,
        MatTableModule,
        MatPaginatorModule,
        MatFormFieldModule,
        MatTooltipModule,
        ReactiveFormsModule,
        MatProgressSpinnerModule,
        FormsModule,
        TranslateModule,
        MatSelectModule,
        DetailLayoutModule,
        HasRolePipeModule,
    ],
})
export class ProjectGrantDetailModule { }
