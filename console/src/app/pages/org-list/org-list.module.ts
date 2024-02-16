import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';
import { TranslateModule } from '@ngx-translate/core';
import { OrgTableModule } from 'src/app/modules/org-table/org-table.module';

import { OrgListRoutingModule } from './org-list-routing.module';
import { OrgListComponent } from './org-list.component';

@NgModule({
  declarations: [OrgListComponent],
  imports: [CommonModule, OrgListRoutingModule, OrgTableModule, TranslateModule],
  exports: [OrgListComponent],
})
export default class OrgListModule {}
