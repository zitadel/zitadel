import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { MatIconModule } from '@angular/material/icon';
import { MatLegacyButtonModule as MatButtonModule } from '@angular/material/legacy-button';
import { MatLegacyCheckboxModule as MatCheckboxModule } from '@angular/material/legacy-checkbox';
import { MatLegacyChipsModule as MatChipsModule } from '@angular/material/legacy-chips';
import { MatLegacyProgressSpinnerModule as MatProgressSpinnerModule } from '@angular/material/legacy-progress-spinner';
import { MatLegacySelectModule as MatSelectModule } from '@angular/material/legacy-select';
import { MatLegacyTableModule as MatTableModule } from '@angular/material/legacy-table';
import { MatLegacyTooltipModule as MatTooltipModule } from '@angular/material/legacy-tooltip';
import { MatSortModule } from '@angular/material/sort';
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
