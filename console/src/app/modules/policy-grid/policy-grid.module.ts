import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';
import { MatButtonModule } from '@angular/material/button';
import { MatIconModule } from '@angular/material/icon';
import { RouterModule } from '@angular/router';
import { TranslateModule } from '@ngx-translate/core';
import { HasRoleModule } from 'src/app/directives/has-role/has-role.module';
import { HasRolePipeModule } from 'src/app/pipes/has-role-pipe/has-role-pipe.module';

import { PolicyGridComponent } from './policy-grid.component';

@NgModule({
    declarations: [PolicyGridComponent],
    imports: [
        CommonModule,
        HasRolePipeModule,
        HasRoleModule,
        TranslateModule,
        RouterModule,
        MatButtonModule,
        MatIconModule,
    ],
    exports: [
        PolicyGridComponent,
    ],
})
export class PolicyGridModule { }
