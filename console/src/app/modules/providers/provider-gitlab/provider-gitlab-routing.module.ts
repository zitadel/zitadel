import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';
import { ProviderGitlabComponent } from './provider-gitlab.component';

const routes: Routes = [
  {
    path: '',
    component: ProviderGitlabComponent,
    data: { animation: 'DetailPage' },
  },
];

@NgModule({
  imports: [RouterModule.forChild(routes)],
  exports: [RouterModule],
})
export class ProviderGitlabRoutingModule {}
