import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';
import { TranslateModule } from '@ngx-translate/core';
import { OrgTableModule } from 'src/app/modules/org-table/org-table.module';

import { ActionsRoutingModule } from './actions-routing.module';
import { ActionsComponent } from './actions.component';
import { MetaLayoutModule } from 'src/app/modules/meta-layout/meta-layout.module';
import { SidenavModule } from 'src/app/modules/sidenav/sidenav.module';

@NgModule({
  declarations: [ActionsComponent],
  imports: [CommonModule, ActionsRoutingModule, OrgTableModule, TranslateModule, MetaLayoutModule, SidenavModule],
  exports: [ActionsComponent],
})
export default class ActionsModule {}
