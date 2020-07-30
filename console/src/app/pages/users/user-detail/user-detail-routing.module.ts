import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';
import { AuthGuard } from 'src/app/guards/auth.guard';
import { RoleGuard } from 'src/app/guards/role.guard';

import { AuthUserDetailComponent } from './auth-user-detail/auth-user-detail.component';
import { PasswordComponent } from './password/password.component';
import { UserDetailComponent } from './user-detail/user-detail.component';

const routes: Routes = [
    {
        path: 'create',
        loadChildren: () => import('../user-create/user-create.module').then(m => m.UserCreateModule),
        canActivate: [AuthGuard, RoleGuard],
        data: {
            roles: ['user.write'],
        },
    },
    {
        path: 'me',
        component: AuthUserDetailComponent,
        canActivate: [AuthGuard],
        data: {
            animation: 'HomePage',
        },
    },
    {
        path: 'me/password',
        component: PasswordComponent,
        canActivate: [AuthGuard],
        data: { animation: 'AddPage' },
    },
    {
        path: ':id',
        component: UserDetailComponent,
        canActivate: [AuthGuard, RoleGuard],
        data: {
            roles: ['user.read'],
            animation: 'HomePage',
        },
    },
    {
        path: ':id/password',
        component: PasswordComponent,
        canActivate: [AuthGuard, RoleGuard],
        data: {
            roles: ['user.write'],
            animation: 'AddPage',
        },
    },
];

@NgModule({
    imports: [RouterModule.forChild(routes)],
    exports: [RouterModule],
})
export class UserDetailRoutingModule { }
