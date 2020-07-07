import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';
import { RoleGuard } from 'src/app/guards/role.guard';

import { UserDetailComponent } from './user-detail/user-detail.component';

const routes: Routes = [
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
        loadChildren: () => import('src/app/pages/password/password.module').then(m => m.PasswordModule),
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
