import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { MatButtonModule } from '@angular/material/button';
import { MatIconModule } from '@angular/material/icon';
import { MatRadioModule } from '@angular/material/radio';
import { MatSortModule } from '@angular/material/sort';
import { MatTableModule } from '@angular/material/table';
import { MatTooltipModule } from '@angular/material/tooltip';
import { TranslateModule } from '@ngx-translate/core';
import { CopyToClipboardModule } from 'src/app/directives/copy-to-clipboard/copy-to-clipboard.module';
import { HasRoleModule } from 'src/app/directives/has-role/has-role.module';
import { OrgListRoutingModule } from 'src/app/pages/orgs/org-list/org-list-routing.module';
import { LocalizedDatePipeModule } from 'src/app/pipes/localized-date-pipe/localized-date-pipe.module';
import { TimestampToDatePipeModule } from 'src/app/pipes/timestamp-to-date-pipe/timestamp-to-date-pipe.module';

import { ActionKeysModule } from '../action-keys/action-keys.module';
import { FilterOrgModule } from '../filter-org/filter-org.module';
import { InputModule } from '../input/input.module';
import { PaginatorModule } from '../paginator/paginator.module';
import { RefreshTableModule } from '../refresh-table/refresh-table.module';
import { OrgsComponent } from './orgs.component';

@NgModule({
  declarations: [OrgsComponent],
  imports: [
    CommonModule,
    OrgListRoutingModule,
    MatTableModule,
    TranslateModule,
    RefreshTableModule,
    ActionKeysModule,
    FilterOrgModule,
    TimestampToDatePipeModule,
    LocalizedDatePipeModule,
    MatSortModule,
    MatIconModule,
    PaginatorModule,
    HasRoleModule,
    MatButtonModule,
    MatTooltipModule,
    CopyToClipboardModule,
    MatRadioModule,
    InputModule,
    FormsModule,
  ],
  exports: [OrgsComponent],
})
export class OrgsModule {}
