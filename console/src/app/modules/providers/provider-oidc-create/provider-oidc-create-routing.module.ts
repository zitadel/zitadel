import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';

import { ProviderOIDCCreateComponent } from './provider-oidc-create.component';

const routes: Routes = [
  {
    path: '',
    component: ProviderOIDCCreateComponent,
    data: { animation: 'DetailPage' },
  },
];

@NgModule({
  imports: [RouterModule.forChild(routes)],
  exports: [RouterModule],
})
export class ProviderOIDCCreateRoutingModule {}
