import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';
import { ProviderGithubESComponent } from './provider-github-es.component';

const routes: Routes = [
  {
    path: '',
    component: ProviderGithubESComponent,
    data: { animation: 'DetailPage' },
  },
];

@NgModule({
  imports: [RouterModule.forChild(routes)],
  exports: [RouterModule],
})
export class ProviderGithubESRoutingModule {}
