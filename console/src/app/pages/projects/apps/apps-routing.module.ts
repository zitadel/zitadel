import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';

import { AppDetailComponent } from '../apps/app-detail/app-detail.component';
import { IntegrateAppComponent } from './integrate/integrate.component';
import { AppCreateComponent } from '../apps/app-create/app-create.component';

const routes: Routes = [
  {
    path: 'integrate',
    component: IntegrateAppComponent,
    data: { animation: 'AddPage' },
  },
  {
    path: 'create',
    component: AppCreateComponent,
    data: { animation: 'AddPage' },
  },
  {
    path: ':appid',
    component: AppDetailComponent,
    data: { animation: 'HomePage' },
  },
];

@NgModule({
  imports: [RouterModule.forChild(routes)],
  exports: [RouterModule],
})
export class AppsRoutingModule {}
