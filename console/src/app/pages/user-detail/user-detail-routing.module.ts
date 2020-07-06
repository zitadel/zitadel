import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';
import { RoleGuard } from 'src/app/guards/role.guard';

import { AuthUserDetailComponent } from './auth-user-detail/auth-user-detail.component';
import { PasswordComponent } from './password/password.component';
import { UserDetailComponent } from './user-detail/user-detail.component';

const routes: Routes = [
    {
        path: 'me',
        component: AuthUserDetailComponent,
    },
    {
        path: 'me/password',
        component: PasswordComponent,
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
        path: ':id/password',
        component: PasswordComponent,
        canActivate: [RoleGuard],
        data: {
            roles: ['user.write'],
        },
    },
];

@NgModule({
    imports: [RouterModule.forChild(routes)],
    exports: [RouterModule],
})
export class UserDetailRoutingModule { }
