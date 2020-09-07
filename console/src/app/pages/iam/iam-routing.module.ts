import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';
import { AuthGuard } from 'src/app/guards/auth.guard';
import { RoleGuard } from 'src/app/guards/role.guard';

import { IamComponent } from './iam.component';

const routes: Routes = [
    {
        path: '',
        component: IamComponent,
        canActivate: [AuthGuard, RoleGuard],
        data: {
            roles: ['iam.read'],
        },
    },
    {
        path: 'members',
        loadChildren: () => import('./iam-members/iam-members.module').then(m => m.IamMembersModule),
        canActivate: [AuthGuard, RoleGuard],
        data: {
            roles: ['iam.member.read'],
        },
    },
    {
        path: 'idp/create',
        loadChildren: () => import('./idp-create/idp-create.module').then(m => m.IdpCreateModule),
        canActivate: [AuthGuard, RoleGuard],
        data: {
            roles: ['iam.write'],
        },
    },
];

@NgModule({
    imports: [RouterModule.forChild(routes)],
    exports: [RouterModule],
})
export class IamRoutingModule { }
