import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';
import { FormsModule, ReactiveFormsModule } from '@angular/forms';
import { MatAutocompleteModule } from '@angular/material/autocomplete';
import { MatButtonModule } from '@angular/material/button';
import { MatCheckboxModule } from '@angular/material/checkbox';
import { MatChipsModule } from '@angular/material/chips';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatIconModule } from '@angular/material/icon';
import { MatProgressSpinnerModule } from '@angular/material/progress-spinner';
import { MatSelectModule } from '@angular/material/select';
import { MatTooltipModule } from '@angular/material/tooltip';
import { TranslateModule } from '@ngx-translate/core';
import { HasRoleModule } from 'src/app/directives/has-role/has-role.module';
import { DetailLayoutModule } from 'src/app/modules/detail-layout/detail-layout.module';
import { MembersTableModule } from 'src/app/modules/members-table/members-table.module';
import { HasRolePipeModule } from 'src/app/pipes/has-role-pipe/has-role-pipe.module';

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
        MatTooltipModule,
        ReactiveFormsModule,
        MatProgressSpinnerModule,
        FormsModule,
        TranslateModule,
        MatFormFieldModule,
        MatSelectModule,
        HasRolePipeModule,
        MembersTableModule,
    ],
})
export class IamMembersModule { }
