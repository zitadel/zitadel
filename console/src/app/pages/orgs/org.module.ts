import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';
import { FormsModule, ReactiveFormsModule } from '@angular/forms';
import { MatButtonModule } from '@angular/material/button';
import { MatDialogModule } from '@angular/material/dialog';
import { MatIconModule } from '@angular/material/icon';
import { MatMenuModule } from '@angular/material/menu';
import { MatProgressSpinnerModule } from '@angular/material/progress-spinner';
import { MatTooltipModule } from '@angular/material/tooltip';
import { TranslateModule } from '@ngx-translate/core';
import { CopyToClipboardModule } from 'src/app/directives/copy-to-clipboard/copy-to-clipboard.module';
import { HasRoleModule } from 'src/app/directives/has-role/has-role.module';
import { MemberCreateDialogModule } from 'src/app/modules/add-member-dialog/member-create-dialog.module';
import { CardModule } from 'src/app/modules/card/card.module';
import { ContributorsModule } from 'src/app/modules/contributors/contributors.module';
import { InfoRowModule } from 'src/app/modules/info-row/info-row.module';
import { InfoSectionModule } from 'src/app/modules/info-section/info-section.module';
import { InputModule } from 'src/app/modules/input/input.module';
import { MetaLayoutModule } from 'src/app/modules/meta-layout/meta-layout.module';
import { MetadataModule } from 'src/app/modules/metadata/metadata.module';
import { NameDialogModule } from 'src/app/modules/name-dialog/name-dialog.module';
import { SettingsGridModule } from 'src/app/modules/settings-grid/settings-grid.module';
import { TopViewModule } from 'src/app/modules/top-view/top-view.module';
import { WarnDialogModule } from 'src/app/modules/warn-dialog/warn-dialog.module';
import { HasRolePipeModule } from 'src/app/pipes/has-role-pipe/has-role-pipe.module';

import { ChangesModule } from '../../modules/changes/changes.module';
import { OrgDetailComponent } from './org-detail/org-detail.component';
import { OrgRoutingModule } from './org-routing.module';

@NgModule({
  declarations: [OrgDetailComponent],
  imports: [
    CommonModule,
    HasRolePipeModule,
    OrgRoutingModule,
    FormsModule,
    InfoRowModule,
    HasRoleModule,
    InputModule,
    InfoSectionModule,
    MatButtonModule,
    MatDialogModule,
    CardModule,
    TopViewModule,
    MatIconModule,
    ReactiveFormsModule,
    MetaLayoutModule,
    MatTooltipModule,
    WarnDialogModule,
    MemberCreateDialogModule,
    MatMenuModule,
    NameDialogModule,
    ChangesModule,
    MatProgressSpinnerModule,
    MetadataModule,
    TranslateModule,
    SettingsGridModule,
    ContributorsModule,
    CopyToClipboardModule,
  ],
})
export default class OrgModule {}
