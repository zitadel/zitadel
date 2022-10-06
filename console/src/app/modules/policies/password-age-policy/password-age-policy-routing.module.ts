import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';

import { PasswordAgePolicyComponent } from './password-age-policy.component';

const routes: Routes = [
  {
    path: '',
    component: PasswordAgePolicyComponent,
    data: {
      animation: 'DetailPage',
    },
  },
];

@NgModule({
  imports: [RouterModule.forChild(routes)],
  exports: [RouterModule],
})
export class PasswordAgePolicyRoutingModule {}
