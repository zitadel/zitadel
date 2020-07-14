import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';
import { RoleGuard } from 'src/app/guards/role.guard';

const routes: Routes = [
    {
        path: 'all',
        loadChildren: () => import('src/app/pages/users/user-list/user-list.module').then(m => m.UserListModule),
        canActivate: [RoleGuard],
        data: {
            roles: ['user.read'],
        },
    },
    {
        path: '',
        loadChildren: () => import('src/app/pages/users/user-detail/user-detail.module').then(m => m.UserDetailModule),
    },
];

@NgModule({
    imports: [RouterModule.forChild(routes)],
    exports: [RouterModule],
})
export class UsersRoutingModule { }
