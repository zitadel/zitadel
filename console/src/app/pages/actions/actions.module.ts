import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';
import { TranslateModule } from '@ngx-translate/core';
import { OrgTableModule } from 'src/app/modules/org-table/org-table.module';

import { ActionsRoutingModule } from './actions-routing.module';
import { ActionsComponent } from './actions.component';
import { MetaLayoutModule } from 'src/app/modules/meta-layout/meta-layout.module';
import { SidenavModule } from 'src/app/modules/sidenav/sidenav.module';
import { ActionsTwoActionsComponent } from 'src/app/modules/actions-two/actions-two-actions/actions-two-actions.component';
import ActionsTwoModule from 'src/app/modules/actions-two/actions-two.module';
import { FormsModule } from '@angular/forms';
import { InfoSectionModule } from 'src/app/modules/info-section/info-section.module';

@NgModule({
  declarations: [ActionsComponent],
  imports: [
    CommonModule,
    FormsModule,
    ActionsRoutingModule,
    OrgTableModule,
    TranslateModule,
    InfoSectionModule,
    SidenavModule,
    ActionsTwoModule,
  ],
  exports: [ActionsComponent],
})
export default class ActionsModule {}
