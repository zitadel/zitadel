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
        loadChildren: () => import('./pages/projects/granted-projects/granted-projects.module')
            .then(m => m.GrantedProjectsModule),
        canActivate: [AuthGuard, RoleGuard],
        data: {
            roles: ['project.read'],
        },
    },
    {
        path: 'projects',
        loadChildren: () => import('./pages/projects/owned-projects/owned-projects.module')
            .then(m => m.OwnedProjectsModule),
        canActivate: [AuthGuard, RoleGuard],
        data: {
            roles: ['project.read'],
        },
    },
    {
        path: 'users',
        loadChildren: () => import('./pages/users/users.module').then(m => m.UsersModule),
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
        path: 'grant-create/project/:projectid/grant/:grantid',
        loadChildren: () => import('src/app/pages/user-grant-create/user-grant-create.module')
            .then(m => m.UserGrantCreateModule),
    },
    {
        path: 'grant-create/project/:projectid',
        loadChildren: () => import('src/app/pages/user-grant-create/user-grant-create.module')
            .then(m => m.UserGrantCreateModule),
    },
    {
        path: 'grant-create/user/:userid',
        loadChildren: () => import('src/app/pages/user-grant-create/user-grant-create.module')
            .then(m => m.UserGrantCreateModule),
    },
    {
        path: 'grant-create',
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
