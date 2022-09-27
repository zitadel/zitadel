import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';

import { FailedEventsComponent } from './failed-events.component';

const routes: Routes = [
  {
    path: '',
    component: FailedEventsComponent,
  },
];

@NgModule({
  imports: [RouterModule.forChild(routes)],
  exports: [RouterModule],
})
export class FailedEventsRoutingModule {}
