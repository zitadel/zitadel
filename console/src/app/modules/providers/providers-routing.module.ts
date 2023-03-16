import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';

import { ProvidersComponent } from './providers.component';

const routes: Routes = [
  {
    path: 'create',
    component: ProvidersComponent,
  },
  {
    path: ':id',
    component: ProvidersComponent,
  },
];

@NgModule({
  imports: [RouterModule.forChild(routes)],
  exports: [RouterModule],
})
export class ProvidersRoutingModule {}
