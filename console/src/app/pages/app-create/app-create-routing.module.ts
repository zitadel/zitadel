import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';

import { AppCreateComponent } from './app-create.component';

const routes: Routes = [
  {
    path: '',
    component: AppCreateComponent,
    data: { animation: 'DetailPage' },
  },
];

@NgModule({
  imports: [RouterModule.forChild(routes)],
  exports: [RouterModule],
})
export class AppCreateRoutingModule {}
