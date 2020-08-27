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
import { MatSortModule } from '@angular/material/sort';
import { MatTableModule } from '@angular/material/table';
import { MatTooltipModule } from '@angular/material/tooltip';
import { TranslateModule } from '@ngx-translate/core';
import { HasRoleModule } from 'src/app/directives/has-role/has-role.module';
import { DetailLayoutModule } from 'src/app/modules/detail-layout/detail-layout.module';
import { RefreshTableModule } from 'src/app/modules/refresh-table/refresh-table.module';
import { HasRolePipeModule } from 'src/app/pipes/has-role-pipe.module';

import { IamMembersRoutingModule } from './iam-members-routing.module';
import { IamMembersComponent } from './iam-members.component';


@NgModule({
    declarations: [IamMembersComponent],
    imports: [
        IamMembersRoutingModule,
        DetailLayoutModule,
        CommonModule,
        MatAutocompleteModule,
        MatChipsModule,
        HasRoleModule,
        MatButtonModule,
        MatCheckboxModule,
        MatIconModule,
        MatTableModule,
        MatPaginatorModule,
        MatSortModule,
        MatTooltipModule,
        ReactiveFormsModule,
        MatProgressSpinnerModule,
        FormsModule,
        TranslateModule,
        MatFormFieldModule,
        MatSelectModule,
        HasRolePipeModule,
        RefreshTableModule,
    ],
})
export class IamMembersModule { }
