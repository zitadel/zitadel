import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';

import { OrgCreateComponent } from './org-create.component';

const routes: Routes = [
  {
    path: '',
    component: OrgCreateComponent,
    data: { animation: 'DetailPage' },
  },
];

@NgModule({
  imports: [RouterModule.forChild(routes)],
  exports: [RouterModule],
})
export class OrgCreateRoutingModule {}
