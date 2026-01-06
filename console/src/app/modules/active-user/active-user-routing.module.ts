import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';
import { authGuard } from 'src/app/guards/auth.guard';
import { roleGuard } from 'src/app/guards/role-guard';
import { ActiveUserPageComponent } from './active-user-page.component';

const routes: Routes = [
  {
    path: '',
    component: ActiveUserPageComponent,
    canActivate: [authGuard, roleGuard],
    data: {
      roles: ['org.read'],
    },
  },
];

@NgModule({
  imports: [RouterModule.forChild(routes)],
  exports: [RouterModule],
})
export default class ActiveUserRoutingModule {}
