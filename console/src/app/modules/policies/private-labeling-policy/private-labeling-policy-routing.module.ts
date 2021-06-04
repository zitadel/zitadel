import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';

import { PrivateLabelingPolicyComponent } from './private-labeling-policy.component';

const routes: Routes = [
  {
    path: '',
    component: PrivateLabelingPolicyComponent,
    data: {
      animation: 'DetailPage',
    },
  },
];

@NgModule({
  imports: [RouterModule.forChild(routes)],
  exports: [RouterModule],
})
export class PrivateLabelingPolicyRoutingModule { }
