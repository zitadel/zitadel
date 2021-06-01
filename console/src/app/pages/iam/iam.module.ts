import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';
import { FormsModule, ReactiveFormsModule } from '@angular/forms';
import { MatAutocompleteModule } from '@angular/material/autocomplete';
import { MatButtonModule } from '@angular/material/button';
import { MatCheckboxModule } from '@angular/material/checkbox';
import { MatChipsModule } from '@angular/material/chips';
import { MatDialogModule } from '@angular/material/dialog';
import { MatIconModule } from '@angular/material/icon';
import { MatPaginatorModule } from '@angular/material/paginator';
import { MatProgressSpinnerModule } from '@angular/material/progress-spinner';
import { MatSortModule } from '@angular/material/sort';
import { MatTableModule } from '@angular/material/table';
import { MatTooltipModule } from '@angular/material/tooltip';
import { TranslateModule } from '@ngx-translate/core';
import { HasRoleModule } from 'src/app/directives/has-role/has-role.module';
import { CardModule } from 'src/app/modules/card/card.module';
import { ChangesModule } from 'src/app/modules/changes/changes.module';
import { ContributorsModule } from 'src/app/modules/contributors/contributors.module';
import { InputModule } from 'src/app/modules/input/input.module';
import { MetaLayoutModule } from 'src/app/modules/meta-layout/meta-layout.module';
import { PolicyGridModule } from 'src/app/modules/policy-grid/policy-grid.module';
import { RefreshTableModule } from 'src/app/modules/refresh-table/refresh-table.module';
import { SharedModule } from 'src/app/modules/shared/shared.module';
import { ZitadelTierModule } from 'src/app/modules/zitadel-tier/zitadel-tier.module';
import { HasRolePipeModule } from 'src/app/pipes/has-role-pipe/has-role-pipe.module';
import { LocalizedDatePipeModule } from 'src/app/pipes/localized-date-pipe/localized-date-pipe.module';
import { TimestampToDatePipeModule } from 'src/app/pipes/timestamp-to-date-pipe/timestamp-to-date-pipe.module';

import { EventstoreComponent } from './eventstore/eventstore.component';
import { FailedEventsComponent } from './failed-events/failed-events.component';
import { IamRoutingModule } from './iam-routing.module';
import { IamViewsComponent } from './iam-views/iam-views.component';
import { IamComponent } from './iam.component';

@NgModule({
  declarations: [IamComponent, EventstoreComponent, IamViewsComponent, FailedEventsComponent],
  imports: [
    CommonModule,
    IamRoutingModule,
    ChangesModule,
    CardModule,
    MatAutocompleteModule,
    MatChipsModule,
    MatButtonModule,
    HasRoleModule,
    MatCheckboxModule,
    MetaLayoutModule,
    MatIconModule,
    MatTableModule,
    ZitadelTierModule,
    MatPaginatorModule,
    InputModule,
    MatSortModule,
    MatTooltipModule,
    ReactiveFormsModule,
    MatProgressSpinnerModule,
    FormsModule,
    TranslateModule,
    MatDialogModule,
    ContributorsModule,
    LocalizedDatePipeModule,
    TimestampToDatePipeModule,
    SharedModule,
    RefreshTableModule,
    HasRolePipeModule,
    MatSortModule,
    PolicyGridModule,
  ],
})
export class IamModule { }
