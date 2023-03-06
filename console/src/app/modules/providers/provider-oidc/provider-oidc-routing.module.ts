import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';

import { ProviderOIDCComponent } from './provider-oidc.component';

const routes: Routes = [
  {
    path: '',
    component: ProviderOIDCComponent,
    data: { animation: 'DetailPage' },
  },
];

@NgModule({
  imports: [RouterModule.forChild(routes)],
  exports: [RouterModule],
})
export class ProviderOIDCRoutingModule {}
