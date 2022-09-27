import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';

import { DomainsComponent } from './domains.component';

const routes: Routes = [
  {
    path: '',
    component: DomainsComponent,
  },
];

@NgModule({
  imports: [RouterModule.forChild(routes)],
  exports: [RouterModule],
})
export class DomainsRoutingModule {}
