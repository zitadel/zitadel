import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';

import { IdpCreateComponent } from './idp-create.component';

const routes: Routes = [
  {
    path: '',
    component: IdpCreateComponent,
    data: { animation: 'DetailPage' },
  },
];

@NgModule({
  imports: [RouterModule.forChild(routes)],
  exports: [RouterModule],
})
export class IdpCreateRoutingModule {}
