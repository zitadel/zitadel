import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';
import { ProviderGitlabSelfHostedComponent } from './provider-gitlab-self-hosted.component';

const routes: Routes = [
  {
    path: '',
    component: ProviderGitlabSelfHostedComponent,
    data: { animation: 'DetailPage' },
  },
];

@NgModule({
  imports: [RouterModule.forChild(routes)],
  exports: [RouterModule],
})
export class ProviderGitlabSelfHostedRoutingModule {}
