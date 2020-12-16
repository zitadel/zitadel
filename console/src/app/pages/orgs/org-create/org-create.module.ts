import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';
import { FormsModule, ReactiveFormsModule } from '@angular/forms';
import { MatButtonModule } from '@angular/material/button';
import { MatCheckboxModule } from '@angular/material/checkbox';
import { MatIconModule } from '@angular/material/icon';
import { MatSelectModule } from '@angular/material/select';
import { MatSlideToggleModule } from '@angular/material/slide-toggle';
import { TranslateModule } from '@ngx-translate/core';
import { HasRoleModule } from 'src/app/directives/has-role/has-role.module';
import { InputModule } from 'src/app/modules/input/input.module';
import { PasswordComplexityViewModule } from 'src/app/modules/password-complexity-view/password-complexity-view.module';
import { HasRolePipeModule } from 'src/app/pipes/has-role-pipe/has-role-pipe.module';

import { OrgCreateRoutingModule } from './org-create-routing.module';
import { OrgCreateComponent } from './org-create.component';

@NgModule({
    declarations: [OrgCreateComponent],
    imports: [
        OrgCreateRoutingModule,
        CommonModule,
        FormsModule,
        ReactiveFormsModule,
        InputModule,
        MatButtonModule,
        MatIconModule,
        MatSelectModule,
        HasRolePipeModule,
        TranslateModule,
        HasRoleModule,
        MatCheckboxModule,
        PasswordComplexityViewModule,
        MatSlideToggleModule,
    ],
})
export class OrgCreateModule { }
