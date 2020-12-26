import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';
import { AuthGuard } from 'src/app/guards/auth.guard';
import { RoleGuard } from 'src/app/guards/role.guard';
import { ProjectType } from 'src/app/proto/generated/zitadel/management_pb';

import { GrantedProjectDetailComponent } from './granted-project-detail/granted-project-detail.component';
import { GrantedProjectsComponent } from './granted-projects.component';

const routes: Routes = [
    {
        path: '',
        component: GrantedProjectsComponent,
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
        path: ':projectid/grant/:grantid/members',
        data: {
            type: ProjectType.PROJECTTYPE_GRANTED,
            roles: ['project.grant.member.read'],
        },
        loadChildren: () => import('src/app/modules/project-members/project-members.module')
            .then(m => m.ProjectMembersModule),
    },
    {
        path: ':id/grant/:grantId',
        component: GrantedProjectDetailComponent,
        data: { animation: 'HomePage' },
    },
    {
        path: ':projectid/roles/create',
        loadChildren: () => import('../project-role-create/project-role-create.module').then(m => m.ProjectRoleCreateModule),
        canActivate: [RoleGuard],
        data: {
            roles: ['project.write'],
        },
    },
    {
        path: ':projectid/grants/create',
        loadChildren: () => import('../project-grant-create/project-grant-create.module')
            .then(m => m.ProjectGrantCreateModule),
    },
];

@NgModule({
    imports: [RouterModule.forChild(routes)],
    exports: [RouterModule],
})
export class GrantedProjectsRoutingModule { }
