import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';
import { authGuard } from 'src/app/guards/auth.guard';
import { roleGuard } from 'src/app/guards/role-guard';
import { userGuard } from 'src/app/guards/user-guard';
import { Type } from 'src/app/proto/generated/zitadel/user_pb';

import { AuthUserDetailComponent } from './user-detail/auth-user-detail/auth-user-detail.component';
import { PasswordComponent } from './user-detail/password/password.component';
import { UserDetailComponent } from './user-detail/user-detail/user-detail.component';
import { UserListComponent } from './user-list/user-list.component';

const routes: Routes = [
  {
    path: '',
    component: UserListComponent,
    data: {
      animation: 'HomePage',
      type: Type.TYPE_HUMAN,
    },
  },
  {
    path: 'create',
    loadChildren: () => import('./user-create/user-create.module'),
    canActivate: [authGuard, roleGuard],
    data: {
      roles: ['user.write'],
    },
  },
  {
    path: 'create-machine',
    loadChildren: () => import('./user-create-machine/user-create-machine.module'),
    canActivate: [authGuard, roleGuard],
    data: {
      roles: ['user.write'],
    },
  },
  {
    path: 'me',
    component: AuthUserDetailComponent,
    canActivate: [authGuard],
    data: {
      animation: 'HomePage',
    },
  },
  {
    path: 'me/password',
    component: PasswordComponent,
    canActivate: [authGuard],
    data: { animation: 'AddPage' },
  },
  {
    path: ':id',
    component: UserDetailComponent,
    canActivate: [authGuard, userGuard, roleGuard],
    data: {
      roles: ['user.read'],
      animation: 'HomePage',
    },
  },
  {
    path: ':id/password',
    component: PasswordComponent,
    canActivate: [authGuard, roleGuard],
    data: {
      roles: ['user.write'],
      animation: 'AddPage',
    },
  },
];

@NgModule({
  imports: [RouterModule.forChild(routes)],
  exports: [RouterModule],
})
export class UsersRoutingModule {}
