import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';
import { SignalsComponent } from './signals.component';
import { SignalsOverviewComponent } from './signals-overview.component';
import { SignalsQueryComponent } from './signals-query.component';
import { SignalsLogsComponent } from './signals-logs.component';
import { SignalsActivityComponent } from './signals-activity.component';

const routes: Routes = [
  {
    path: '',
    component: SignalsComponent,
    children: [
      {
        path: '',
        component: SignalsOverviewComponent,
      },
      {
        path: 'explore',
        component: SignalsQueryComponent,
      },
      {
        path: 'logs',
        component: SignalsLogsComponent,
      },
      {
        path: 'activity',
        component: SignalsActivityComponent,
      },
    ],
  },
];

@NgModule({
  imports: [RouterModule.forChild(routes)],
  exports: [RouterModule],
})
export class SignalsRoutingModule {}
