import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';
import { RoleGuard } from 'src/app/guards/role.guard';

const routes: Routes = [
    {
        path: '',
        loadChildren: () => import('src/app/pages/user-list/user-list.module').then(m => m.UserListModule),
        canActivate: [RoleGuard],
        data: {
            roles: ['user.read'],
        },
    },
    {
        path: 'me/password',
        loadChildren: () => import('src/app/pages/users/password/password.module').then(m => m.PasswordModule),
        data: {
            roles: ['user.write'],
        },
    },
    {
        path: 'me',
        loadChildren: () => import('src/app/pages/users/auth-user-detail/auth-user-detail.module')
            .then(m => m.AuthUserDetailModule),
        data: {
            roles: ['user.write'],
        },
    },
    {
        path: ':id/password',
        loadChildren: () => import('src/app/pages/users/password/password.module').then(m => m.PasswordModule),
        data: {
            roles: ['user.write'],
        },
    },
    {
        path: ':id',
        loadChildren: () => import('src/app/pages/users/auth-user-detail/auth-user-detail.module')
            .then(m => m.AuthUserDetailModule),
        canActivate: [RoleGuard],
        data: {
            roles: ['user.read'],
        },
    },
];

@NgModule({
    imports: [RouterModule.forChild(routes)],
    exports: [RouterModule],
})
export class UsersRoutingModule { }
