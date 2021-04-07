import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';
import { MatButtonModule } from '@angular/material/button';
import { MatCheckboxModule } from '@angular/material/checkbox';
import { MatIconModule } from '@angular/material/icon';
import { MatPaginatorModule } from '@angular/material/paginator';
import { MatProgressSpinnerModule } from '@angular/material/progress-spinner';
import { MatTableModule } from '@angular/material/table';
import { MatTooltipModule } from '@angular/material/tooltip';
import { RouterModule, Routes } from '@angular/router';
import { TranslateModule } from '@ngx-translate/core';
import { HasRoleModule } from 'src/app/directives/has-role/has-role.module';
import { DetailLayoutModule } from 'src/app/modules/detail-layout/detail-layout.module';
import { RefreshTableModule } from 'src/app/modules/refresh-table/refresh-table.module';
import { HasRolePipeModule } from 'src/app/pipes/has-role-pipe/has-role-pipe.module';
import { LocalizedDatePipeModule } from 'src/app/pipes/localized-date-pipe/localized-date-pipe.module';
import { TimestampToDatePipeModule } from 'src/app/pipes/timestamp-to-date-pipe/timestamp-to-date-pipe.module';

import { MembershipDetailComponent } from './membership-detail.component';

const routes: Routes = [
    {
        path: '',
        component: MembershipDetailComponent,
        canActivate: [],
        data: {
            roles: ['user.write'],
        },
    },
];
@NgModule({
    declarations: [MembershipDetailComponent],
    imports: [
        CommonModule,
        RouterModule.forChild(routes),
        TranslateModule,
        DetailLayoutModule,
        MatCheckboxModule,
        MatTableModule,
        MatPaginatorModule,
        MatProgressSpinnerModule,
        LocalizedDatePipeModule,
        TimestampToDatePipeModule,
        HasRoleModule,
        MatIconModule,
        MatButtonModule,
        HasRolePipeModule,
        RefreshTableModule,
        MatTooltipModule,
    ],
})
export class MembershipDetailModule { }
