import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { MatButtonModule } from '@angular/material/button';
import { MatCheckboxModule } from '@angular/material/checkbox';
import { MatChipsModule } from '@angular/material/chips';
import { MatIconModule } from '@angular/material/icon';
import { MatProgressSpinnerModule } from '@angular/material/progress-spinner';
import { MatSelectModule } from '@angular/material/select';
import { MatSortModule } from '@angular/material/sort';
import { MatTableModule } from '@angular/material/table';
import { MatTooltipModule } from '@angular/material/tooltip';
import { RouterModule } from '@angular/router';
import { TranslateModule } from '@ngx-translate/core';
import { InputModule } from 'src/app/modules/input/input.module';
import { HasRolePipeModule } from 'src/app/pipes/has-role-pipe/has-role-pipe.module';
import { RoleTransformPipeModule } from 'src/app/pipes/role-transform/role-transform.module';

import { AvatarModule } from '../avatar/avatar.module';
import { PaginatorModule } from '../paginator/paginator.module';
import { RefreshTableModule } from '../refresh-table/refresh-table.module';
import { MembershipsTableComponent } from './memberships-table.component';

@NgModule({
  declarations: [MembershipsTableComponent],
  imports: [
    CommonModule,
    InputModule,
    MatSelectModule,
    MatCheckboxModule,
    MatIconModule,
    PaginatorModule,
    MatChipsModule,
    MatTooltipModule,
    RoleTransformPipeModule,
    FormsModule,
    TranslateModule,
    RefreshTableModule,
    RouterModule,
    AvatarModule,
    MatTableModule,
    MatSortModule,
    MatProgressSpinnerModule,
    MatButtonModule,
    HasRolePipeModule,
  ],
  exports: [MembershipsTableComponent],
})
export class MembershipsTableModule {}
