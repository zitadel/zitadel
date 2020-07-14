import { CommonModule } from '@angular/common';
import { NgModule, NO_ERRORS_SCHEMA } from '@angular/core';
import { MatButtonModule } from '@angular/material/button';
import { MatIconModule } from '@angular/material/icon';
import { MatTabsModule } from '@angular/material/tabs';
import { TranslateModule } from '@ngx-translate/core';
import { HasRoleModule } from 'src/app/directives/has-role/has-role.module';
import { MetaLayoutModule } from 'src/app/modules/meta-layout/meta-layout.module';
import { HasRolePipeModule } from 'src/app/pipes/has-role-pipe.module';

import { ApplicationGridComponent } from './application-grid/application-grid.component';
import { ApplicationsComponent } from './applications/applications.component';
import { OwnedProjectDetailRoutingModule } from './owned-project-detail-routing.module';
import { OwnedProjectDetailComponent } from './owned-project-detail.component';

@NgModule({
    declarations: [
        OwnedProjectDetailComponent,
        ApplicationGridComponent,
        ApplicationsComponent,
    ],
    imports: [
        CommonModule,
        OwnedProjectDetailRoutingModule,
        TranslateModule,
        HasRolePipeModule,
        HasRoleModule,
        MatTabsModule,
        MatButtonModule,
        MatIconModule,
        MetaLayoutModule,
    ],
    schemas: [NO_ERRORS_SCHEMA],
})
export class OwnedProjectDetailModule { }
