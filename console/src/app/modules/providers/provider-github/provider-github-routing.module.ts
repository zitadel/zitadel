import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';

import { ProviderGithubComponent } from './provider-github.component';

const routes: Routes = [
  {
    path: '',
    component: ProviderGithubComponent,
    data: { animation: 'DetailPage' },
  },
];

@NgModule({
  imports: [RouterModule.forChild(routes)],
  exports: [RouterModule],
})
export class ProviderGithubRoutingModule {}
