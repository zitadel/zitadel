import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';

import { ProviderGoogleComponent } from './provider-google.component';

const routes: Routes = [
  {
    path: '',
    component: ProviderGoogleComponent,
    data: { animation: 'DetailPage' },
  },
];

@NgModule({
  imports: [RouterModule.forChild(routes)],
  exports: [RouterModule],
})
export class ProviderGoogleRoutingModule {}
