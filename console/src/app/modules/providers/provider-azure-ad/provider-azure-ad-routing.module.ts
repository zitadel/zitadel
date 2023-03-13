import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';

import { ProviderAzureADComponent } from './provider-azure-ad.component';

const routes: Routes = [
  {
    path: '',
    component: ProviderAzureADComponent,
    data: { animation: 'DetailPage' },
  },
];

@NgModule({
  imports: [RouterModule.forChild(routes)],
  exports: [RouterModule],
})
export class ProviderAzureADRoutingModule {}
