import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';
import { AuthGuard } from 'src/app/guards/auth.guard';
import { RoleGuard } from 'src/app/guards/role.guard';
import { ProjectType } from 'src/app/proto/generated/management_pb';

import { OwnedProjectDetailComponent } from './owned-project-detail/owned-project-detail.component';
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
        canActivate: [AuthGuard, RoleGuard],
        data: {
            roles: ['project.write'],
        },
    },
    {
        path: ':id',
        component: OwnedProjectDetailComponent,
        data: { animation: 'HomePage' },
    },
    {
        path: ':projectid/members',
        data: {
            type: ProjectType.PROJECTTYPE_OWNED,
        },
        loadChildren: () => import('../../modules/project-members/project-members.module').then(m => m.ProjectMembersModule),
    },
    {
        path: ':projectid/apps',
        data: { animation: 'AddPage' },
        loadChildren: () => import('../apps/apps.module').then(m => m.AppsModule),
    },
    {
        path: ':projectid/roles/create',
        loadChildren: () => import('../project-role-create/project-role-create.module').then(m => m.ProjectRoleCreateModule),

    },
    {
        path: ':projectid/grants/create',
        loadChildren: () => import('../project-grant-create/project-grant-create.module')
            .then(m => m.ProjectGrantCreateModule),
    },
    {
        path: ':projectid/grant/:grantid',
        loadChildren: () => import('./project-grant-detail/project-grant-detail.module')
            .then(m => m.ProjectGrantDetailModule),
    },
];

@NgModule({
    imports: [RouterModule.forChild(routes)],
    exports: [RouterModule],
})
export class OwnedProjectsRoutingModule { }
