import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';
import { FormsModule, ReactiveFormsModule } from '@angular/forms';
import { MatButtonModule } from '@angular/material/button';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatIconModule } from '@angular/material/icon';
import { MatInputModule } from '@angular/material/input';
import { MatSelectModule } from '@angular/material/select';
import { TranslateModule } from '@ngx-translate/core';
import { HasRolePipeModule } from 'src/app/pipes/has-role-pipe.module';

import { OrgCreateRoutingModule } from './org-create-routing.module';
import { OrgCreateComponent } from './org-create.component';

@NgModule({
    declarations: [OrgCreateComponent],
    imports: [
        OrgCreateRoutingModule,
        CommonModule,
        FormsModule,
        ReactiveFormsModule,
        MatInputModule,
        MatFormFieldModule,
        MatButtonModule,
        MatIconModule,
        MatSelectModule,
        HasRolePipeModule,
        TranslateModule,
    ],
})
export class OrgCreateModule { }
