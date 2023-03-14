import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';

import { ProviderOAuthComponent } from './provider-oauth.component';

const routes: Routes = [
  {
    path: '',
    component: ProviderOAuthComponent,
    data: { animation: 'DetailPage' },
  },
];

@NgModule({
  imports: [RouterModule.forChild(routes)],
  exports: [RouterModule],
})
export class ProviderOAuthRoutingModule {}
