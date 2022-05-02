import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';
import { TranslateModule } from '@ngx-translate/core';
import { OrgsModule } from 'src/app/modules/orgs/orgs.module';

import { OrgListRoutingModule } from './org-list-routing.module';
import { OrgListComponent } from './org-list.component';

@NgModule({
  declarations: [OrgListComponent],
  imports: [CommonModule, OrgListRoutingModule, OrgsModule, TranslateModule],
})
export class OrgListModule {}
