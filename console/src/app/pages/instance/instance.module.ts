import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';
import { FormsModule, ReactiveFormsModule } from '@angular/forms';
import { MatIconModule } from '@angular/material/icon';
import { MatLegacyAutocompleteModule as MatAutocompleteModule } from '@angular/material/legacy-autocomplete';
import { MatLegacyButtonModule as MatButtonModule } from '@angular/material/legacy-button';
import { MatLegacyCheckboxModule as MatCheckboxModule } from '@angular/material/legacy-checkbox';
import { MatLegacyChipsModule as MatChipsModule } from '@angular/material/legacy-chips';
import { MatLegacyDialogModule as MatDialogModule } from '@angular/material/legacy-dialog';
import { MatLegacyProgressSpinnerModule as MatProgressSpinnerModule } from '@angular/material/legacy-progress-spinner';
import { MatLegacyTableModule as MatTableModule } from '@angular/material/legacy-table';
import { MatLegacyTooltipModule as MatTooltipModule } from '@angular/material/legacy-tooltip';
import { MatSortModule } from '@angular/material/sort';
import { TranslateModule } from '@ngx-translate/core';
import { HasRoleModule } from 'src/app/directives/has-role/has-role.module';
import { CardModule } from 'src/app/modules/card/card.module';
import { ChangesModule } from 'src/app/modules/changes/changes.module';
import { ContributorsModule } from 'src/app/modules/contributors/contributors.module';
import { InfoRowModule } from 'src/app/modules/info-row/info-row.module';
import { InputModule } from 'src/app/modules/input/input.module';
import { MetaLayoutModule } from 'src/app/modules/meta-layout/meta-layout.module';
import { OrgTableModule } from 'src/app/modules/org-table/org-table.module';
import { RefreshTableModule } from 'src/app/modules/refresh-table/refresh-table.module';
import { SettingsGridModule } from 'src/app/modules/settings-grid/settings-grid.module';
import { TopViewModule } from 'src/app/modules/top-view/top-view.module';
import { HasRolePipeModule } from 'src/app/pipes/has-role-pipe/has-role-pipe.module';
import { LocalizedDatePipeModule } from 'src/app/pipes/localized-date-pipe/localized-date-pipe.module';
import { TimestampToDatePipeModule } from 'src/app/pipes/timestamp-to-date-pipe/timestamp-to-date-pipe.module';

import { IamRoutingModule } from './instance-routing.module';
import { InstanceComponent } from './instance.component';

@NgModule({
  declarations: [InstanceComponent],
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
    TopViewModule,
    MatTableModule,
    InfoRowModule,
    InputModule,
    MatSortModule,
    MatTooltipModule,
    ReactiveFormsModule,
    MatProgressSpinnerModule,
    FormsModule,
    TranslateModule,
    OrgTableModule,
    MatDialogModule,
    ContributorsModule,
    LocalizedDatePipeModule,
    TimestampToDatePipeModule,
    RefreshTableModule,
    HasRolePipeModule,
    MatSortModule,
    SettingsGridModule,
  ],
})
export default class InstanceModule {}
