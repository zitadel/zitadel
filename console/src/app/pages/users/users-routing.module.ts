import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';
import { RoleGuard } from 'src/app/guards/role.guard';

import { AuthGuard } from '../../guards/auth.guard';

const routes: Routes = [
    {
        path: 'all',
        loadChildren: () => import('../../pages/user-list/user-list.module').then(m => m.UserListModule),
        canActivate: [AuthGuard, RoleGuard],
        data: {
            roles: ['user.read'],
        },
    },
    {
        path: '',
        loadChildren: () => import('../user-detail/user-detail.module').then(m => m.UserDetailModule),
        canActivate: [AuthGuard],
    },
];

@NgModule({
    imports: [RouterModule.forChild(routes)],
    exports: [RouterModule],
})
export class UsersRoutingModule { }
