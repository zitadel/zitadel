import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';
import { QuicklinkStrategy } from 'ngx-quicklink';

import { AuthGuard } from './guards/auth.guard';
import { RoleGuard } from './guards/role.guard';

const routes: Routes = [
    {
        path: '',
        loadChildren: () => import('./pages/home/home.module').then(m => m.HomeModule),
        canActivate: [AuthGuard],
    },
    {
        path: 'firststeps',
        loadChildren: () => import('./modules/onboarding/onboarding.module')
            .then(m => m.OnboardingModule),
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
        canActivate: [AuthGuard],
        children: [
            {
                path: 'list',
                loadChildren: () => import('src/app/pages/users/user-list/user-list.module')
                    .then(m => m.UserListModule),
                canActivate: [RoleGuard],
                data: {
                    roles: ['user.read'],
                },
            },
            {
                path: '',
                loadChildren: () => import('src/app/pages/users/user-detail/user-detail.module')
                    .then(m => m.UserDetailModule),
            },
        ],
    },
    {
        path: 'iam',
        loadChildren: () => import('./pages/iam/iam.module').then(m => m.IamModule),
        canActivate: [AuthGuard, RoleGuard],
        data: {
            roles: ['iam.read', 'iam.write'],
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
        path: 'grants',
        loadChildren: () => import('./pages/grants/grants.module').then(m => m.GrantsModule),
        canActivate: [AuthGuard, RoleGuard],
        data: {
            roles: ['user.grant.read'],
        },
    },
    {
        path: 'grant-create',
        canActivate: [AuthGuard],
        children: [
            {
                path: 'project/:projectid/grant/:grantid',
                loadChildren: () => import('src/app/pages/user-grant-create/user-grant-create.module')
                    .then(m => m.UserGrantCreateModule),
                canActivate: [RoleGuard],
                data: {
                    roles: ['user.grant.write'],
                },
            },
            {
                path: 'project/:projectid',
                loadChildren: () => import('src/app/pages/user-grant-create/user-grant-create.module')
                    .then(m => m.UserGrantCreateModule),
                canActivate: [RoleGuard],
                data: {
                    roles: ['user.grant.write'],
                },
            },
            {
                path: 'user/:userid',
                loadChildren: () => import('src/app/pages/user-grant-create/user-grant-create.module')
                    .then(m => m.UserGrantCreateModule),
                canActivate: [RoleGuard],
                data: {
                    roles: ['user.grant.write'],
                },
            },
            {
                path: '',
                loadChildren: () => import('src/app/pages/user-grant-create/user-grant-create.module')
                    .then(m => m.UserGrantCreateModule),
                canActivate: [RoleGuard],
                data: {
                    roles: ['user.grant.write'],
                },
            },
        ],
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
    imports: [
        RouterModule.forRoot(
            routes,
            {
                preloadingStrategy: QuicklinkStrategy,
                relativeLinkResolution: 'legacy',
            },
        ),
    ],
    exports: [RouterModule],
})
export class AppRoutingModule { }
