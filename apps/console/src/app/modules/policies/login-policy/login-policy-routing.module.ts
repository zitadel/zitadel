import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';

import { LoginPolicyComponent } from './login-policy.component';

const routes: Routes = [
  {
    path: '',
    component: LoginPolicyComponent,
    data: {
      animation: 'DetailPage',
    },
  },
];

@NgModule({
  imports: [RouterModule.forChild(routes)],
  exports: [RouterModule],
})
export class LoginPolicyRoutingModule {}
