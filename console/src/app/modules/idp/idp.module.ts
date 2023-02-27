import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';
import { FormsModule, ReactiveFormsModule } from '@angular/forms';
import { MatIconModule } from '@angular/material/icon';
import { MatLegacyButtonModule as MatButtonModule } from '@angular/material/legacy-button';
import { MatLegacyCheckboxModule as MatCheckboxModule } from '@angular/material/legacy-checkbox';
import { MatLegacyChipsModule as MatChipsModule } from '@angular/material/legacy-chips';
import { MatLegacyMenuModule as MatMenuModule } from '@angular/material/legacy-menu';
import { MatLegacySelectModule as MatSelectModule } from '@angular/material/legacy-select';
import { MatLegacyTooltipModule as MatTooltipModule } from '@angular/material/legacy-tooltip';
import { TranslateModule } from '@ngx-translate/core';
import { HasRoleModule } from 'src/app/directives/has-role/has-role.module';
import { DetailLayoutModule } from 'src/app/modules/detail-layout/detail-layout.module';
import { InputModule } from 'src/app/modules/input/input.module';
import { HasRolePipeModule } from 'src/app/pipes/has-role-pipe/has-role-pipe.module';

import { CardModule } from '../card/card.module';
import { InfoRowModule } from '../info-row/info-row.module';
import { InfoSectionModule } from '../info-section/info-section.module';
import { TopViewModule } from '../top-view/top-view.module';
import { WarnDialogModule } from '../warn-dialog/warn-dialog.module';
import { IdpRoutingModule } from './idp-routing.module';
import { IdpComponent } from './idp.component';

@NgModule({
  declarations: [IdpComponent],
  imports: [
    CommonModule,
    IdpRoutingModule,
    FormsModule,
    ReactiveFormsModule,
    InputModule,
    MatButtonModule,
    WarnDialogModule,
    MatIconModule,
    InfoSectionModule,
    MatMenuModule,
    TopViewModule,
    MatTooltipModule,
    MatSelectModule,
    CardModule,
    TranslateModule,
    MatCheckboxModule,
    InfoRowModule,
    MatChipsModule,
    HasRoleModule,
    HasRolePipeModule,
    DetailLayoutModule,
  ],
})
export default class IdpModule {}
