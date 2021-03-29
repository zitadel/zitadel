import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';
import { MatButtonModule } from '@angular/material/button';
import { MatIconModule } from '@angular/material/icon';
import { MatTooltipModule } from '@angular/material/tooltip';
import { RouterModule } from '@angular/router';
import { TranslateModule } from '@ngx-translate/core';
import { HasRoleModule } from 'src/app/directives/has-role/has-role.module';
import { HasFeaturePipeModule } from 'src/app/pipes/has-feature-pipe/has-feature-pipe.module';
import { HasRolePipeModule } from 'src/app/pipes/has-role-pipe/has-role-pipe.module';

import { InfoSectionModule } from '../info-section/info-section.module';
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
        MatTooltipModule,
        InfoSectionModule,
        HasFeaturePipeModule,
    ],
    exports: [
        PolicyGridComponent,
    ],
})
export class PolicyGridModule { }
