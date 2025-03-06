import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';
import { FormsModule, ReactiveFormsModule } from '@angular/forms';
import { MatButtonModule } from '@angular/material/button';
import { MatCheckboxModule } from '@angular/material/checkbox';
import { MatDialogModule } from '@angular/material/dialog';
import { MatIconModule } from '@angular/material/icon';
import { MatMenuModule } from '@angular/material/menu';
import { MatProgressSpinnerModule } from '@angular/material/progress-spinner';
import { MatTableModule } from '@angular/material/table';
import { MatTooltipModule } from '@angular/material/tooltip';
import { RouterModule } from '@angular/router';
import { TranslateModule } from '@ngx-translate/core';
import { CopyToClipboardModule } from 'src/app/directives/copy-to-clipboard/copy-to-clipboard.module';
import { HasRoleModule } from 'src/app/directives/has-role/has-role.module';
import { GroupMemberCreateDialogModule } from 'src/app/modules/add-group-member-dialog/group-member-create-dialog.module';
import { CardModule } from 'src/app/modules/card/card.module';
import { ChangesModule } from 'src/app/modules/changes/changes.module';
import { DetailLayoutModule } from 'src/app/modules/detail-layout/detail-layout.module';
import { InfoSectionModule } from 'src/app/modules/info-section/info-section.module';
import { InputModule } from 'src/app/modules/input/input.module';
import { MetaLayoutModule } from 'src/app/modules/meta-layout/meta-layout.module';
import { PaginatorModule } from 'src/app/modules/paginator/paginator.module';
import { RefreshTableModule } from 'src/app/modules/refresh-table/refresh-table.module';
import { SidenavModule } from 'src/app/modules/sidenav/sidenav.module';
import { TableActionsModule } from 'src/app/modules/table-actions/table-actions.module';
import { TopViewModule } from 'src/app/modules/top-view/top-view.module';
import { GroupGrantsModule } from 'src/app/modules/group-grants/group-grants.module';
import { WarnDialogModule } from 'src/app/modules/warn-dialog/warn-dialog.module';
import { HasRolePipeModule } from 'src/app/pipes/has-role-pipe/has-role-pipe.module';
import { LocalizedDatePipeModule } from 'src/app/pipes/localized-date-pipe/localized-date-pipe.module';
import { TimestampToDatePipeModule } from 'src/app/pipes/timestamp-to-date-pipe/timestamp-to-date-pipe.module';

import { MatSelectModule } from '@angular/material/select';
import { QRCodeModule } from 'angularx-qrcode';
import { InfoDialogModule } from 'src/app/modules/info-dialog/info-dialog.module';
import { MetadataModule } from 'src/app/modules/metadata/metadata.module';
import { ActionKeysModule } from '../../../modules/action-keys/action-keys.module';
import { InfoRowModule } from '../../../modules/info-row/info-row.module';
import { GroupDetailComponent } from './group-detail/group-detail.component';
import { GroupMembersTableModule } from 'src/app/modules/group-members-table/group-members-table.module';


@NgModule({
  declarations: [
    GroupDetailComponent,
  ],
  providers: [],
  imports: [
    ChangesModule,
    CommonModule,
    SidenavModule,
    InfoDialogModule,
    FormsModule,
    ReactiveFormsModule,
    WarnDialogModule,
    MatDialogModule,
    QRCodeModule,
    MetaLayoutModule,
    MatCheckboxModule,
    MetadataModule,
    TopViewModule,
    HasRolePipeModule,
    GroupGrantsModule,
    MatButtonModule,
    MatIconModule,
    CardModule,
    MatProgressSpinnerModule,
    MatTooltipModule,
    HasRoleModule,
    TranslateModule,
    MatTableModule,
    InfoRowModule,
    PaginatorModule,
    MatMenuModule,
    RouterModule,
    RefreshTableModule,
    CopyToClipboardModule,
    DetailLayoutModule,
    TableActionsModule,
    GroupMemberCreateDialogModule,
    TimestampToDatePipeModule,
    LocalizedDatePipeModule,
    InputModule,
    InfoSectionModule,
    MatSelectModule,
    ActionKeysModule,
    GroupMembersTableModule,
  ],
})
export class GroupDetailModule {}
