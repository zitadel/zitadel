import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';

import { AppQuickCreateComponent } from './app-quick-create.component';

const routes: Routes = [
  {
    path: '',
    component: AppQuickCreateComponent,
    data: { animation: 'DetailPage' },
  },
];

@NgModule({
  imports: [RouterModule.forChild(routes)],
  exports: [RouterModule],
})
export class AppCreateRoutingModule {}
