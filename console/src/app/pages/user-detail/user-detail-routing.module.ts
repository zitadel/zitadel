import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';
import { RoleGuard } from 'src/app/guards/role.guard';

import { AuthUserDetailComponent } from './auth-user-detail/auth-user-detail.component';
import { UserDetailComponent } from './user-detail/user-detail.component';

const routes: Routes = [
    {
        path: 'me',
        component: AuthUserDetailComponent,
    },
    {
        path: ':id',
        component: UserDetailComponent,
        canActivate: [RoleGuard],
        data: {
            roles: ['user.read'],
        },
    },
    {
        path: ':id/grant-create',
        loadChildren: () => import('../user-grant-create/user-grant-create.module').then(m => m.UserGrantCreateModule),
    },
    {
        path: ':id/grant/:grantid',
        loadChildren: () => import('./user-grant/user-grant.module').then(m => m.UserGrantModule),
    },
];

@NgModule({
    imports: [RouterModule.forChild(routes)],
    exports: [RouterModule],
})
export class UserDetailRoutingModule { }
