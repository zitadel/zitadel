import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';

import { AuthGuard } from './guards/auth.guard';
import { RoleGuard } from './guards/role.guard';

const routes: Routes = [
    {
        path: '',
        loadChildren: () => import('./pages/home/home.module').then(m => m.HomeModule),
        canActivate: [AuthGuard],
    },
    {
        path: 'granted-projects',
        loadChildren: () => import('./pages/granted-projects/granted-projects.module').then(m => m.GrantedProjectsModule),
        canActivate: [AuthGuard, RoleGuard],
        data: {
            roles: ['project.read'],
        },
    },
    {
        path: 'projects',
        loadChildren: () => import('./pages/owned-projects/owned-projects.module').then(m => m.OwnedProjectsModule),
        canActivate: [AuthGuard, RoleGuard],
        data: {
            roles: ['project.read'],
        },
    },
    {
        path: 'user',
        loadChildren: () => import('./pages/user-detail/user-detail.module').then(m => m.UserDetailModule),
        canActivate: [AuthGuard],
    },
    {
        path: 'users',
        loadChildren: () => import('./pages/user-list/user-list.module').then(m => m.UserListModule),
        canActivate: [AuthGuard, RoleGuard],
        data: {
            roles: ['user.read'],
        },
    },
    {
        path: 'iam',
        loadChildren: () => import('./pages/iam/iam.module').then(m => m.IamModule),
        canActivate: [AuthGuard, RoleGuard],
        data: {
            roles: ['iam.read'],
        },
    },
    {
        path: 'org',
        loadChildren: () => import('./pages/orgs/orgs.module').then(m => m.OrgsModule),
        canActivate: [AuthGuard, RoleGuard],
        data: {
            roles: ['org.read'],
        },
    },
    {
        path: ':id/grant-create',
        loadChildren: () => import('src/app/pages/user-grant-create/user-grant-create.module')
            .then(m => m.UserGrantCreateModule),
    },
    {
        path: 'signedout',
        loadChildren: () => import('./pages/signedout/signedout.module').then(m => m.SignedoutModule),
    },
    {
        path: '**',
        redirectTo: '/',
    },
];

@NgModule({
    imports: [RouterModule.forRoot(routes)],
    exports: [RouterModule],
})
export class AppRoutingModule { }
