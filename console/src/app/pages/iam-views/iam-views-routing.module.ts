import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';

import { IamViewsComponent } from './iam-views.component';

const routes: Routes = [
  {
    path: '',
    component: IamViewsComponent,
  },
];

@NgModule({
  imports: [RouterModule.forChild(routes)],
  exports: [RouterModule],
})
export class IamViewsRoutingModule {}
