import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';
import { authGuard } from 'src/app/guards/auth.guard';
import { roleGuard } from 'src/app/guards/role-guard';

import { GroupDetailComponent } from './group-detail/group-detail/group-detail.component';
import { GroupListComponent } from './group-list/group-list.component';

const routes: Routes = [
  {
    path: '',
    component: GroupListComponent,
    data: {
      animation: 'HomePage',
    },
  },
  {
    path: 'create',
    loadChildren: () => import('./group-create/group-create.module'),
    canActivate: [authGuard, roleGuard],
    data: {
      roles: ['user.write'],
    },
  },
  {
    path: ':id',
    component: GroupDetailComponent,
    canActivate: [authGuard, roleGuard],
    data: {
      roles: ['user.read'],
      animation: 'HomePage',
    },
  },
];

@NgModule({
  imports: [RouterModule.forChild(routes)],
  exports: [RouterModule],
})
export class GroupsRoutingModule { }
