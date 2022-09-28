import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';

import { IdpComponent } from './idp.component';

const routes: Routes = [
  {
    path: '',
    component: IdpComponent,
    data: { animation: 'DetailPage' },
  },
];

@NgModule({
  imports: [RouterModule.forChild(routes)],
  exports: [RouterModule],
})
export class IdpRoutingModule {}
