import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';

import { PrivacyPolicyComponent } from './privacy-policy.component';

const routes: Routes = [
  {
    path: '',
    component: PrivacyPolicyComponent,
    data: {
      animation: 'DetailPage',
    },
  },
];

@NgModule({
  imports: [RouterModule.forChild(routes)],
  exports: [RouterModule],
})
export class PrivacyPolicyRoutingModule { }
