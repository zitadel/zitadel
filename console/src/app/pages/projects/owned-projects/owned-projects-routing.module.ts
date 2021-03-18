import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';
import { RoleGuard } from 'src/app/guards/role.guard';
import { ProjectType } from 'src/app/modules/project-members/project-members-datasource';

import { OwnedProjectsComponent } from './owned-projects.component';

const routes: Routes = [
    {
        path: '',
        component: OwnedProjectsComponent,
        data: { animation: 'HomePage' },
    },
    {
        path: 'create',
        loadChildren: () => import('../project-create/project-create.module').then(m => m.ProjectCreateModule),
        canActivate: [RoleGuard],
        data: {
            roles: ['project.write'],
        },
    },
    {
        path: ':id',
        data: {
            animation: 'HomePage',
            roles: ['project.read'],
        },
        canActivate: [RoleGuard],
        loadChildren: () => import('./owned-project-detail/owned-project-detail.module')
            .then(m => m.OwnedProjectDetailModule),
    },
    {
        path: ':projectid',
        children: [
            {
                path: 'members',
                data: {
                    type: ProjectType.PROJECTTYPE_OWNED,
                    roles: ['project.member.read'],
                },
                canActivate: [RoleGuard],
                loadChildren: () => import('src/app/modules/project-members/project-members.module')
                    .then(m => m.ProjectMembersModule),
            },
            {
                path: 'apps',
                data: {
                    animation: 'AddPage',
                    roles: ['project.app.read'],
                },
                canActivate: [RoleGuard],
                loadChildren: () => import('src/app/pages/projects/apps/apps.module')
                    .then(m => m.AppsModule),
            },
            {
                path: 'roles/create',
                loadChildren: () => import('../project-role-create/project-role-create.module')
                    .then(m => m.ProjectRoleCreateModule),
            },
            {
                path: 'grants/create',
                loadChildren: () => import('../project-grant-create/project-grant-create.module')
                    .then(m => m.ProjectGrantCreateModule),
            },
            {
                path: 'grant/:grantid',
                loadChildren: () => import('./project-grant-detail/project-grant-detail.module')
                    .then(m => m.ProjectGrantDetailModule),
            },
        ],
    },
];

@NgModule({
    imports: [RouterModule.forChild(routes)],
    exports: [RouterModule],
})
export class OwnedProjectsRoutingModule { }
