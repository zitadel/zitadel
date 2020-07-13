import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';
import { AuthGuard } from 'src/app/guards/auth.guard';
import { RoleGuard } from 'src/app/guards/role.guard';

import { UserListComponent } from './user-list.component';


const routes: Routes = [
    {
        path: '',
        component: UserListComponent,
        data: { animation: 'HomePage' },
    },
    {
        path: 'create',
        loadChildren: () => import('../user-create/user-create.module').then(m => m.UserCreateModule),
        canActivate: [AuthGuard, RoleGuard],
        data: {
            roles: ['user.write'],
        },
    },
];

@NgModule({
    imports: [RouterModule.forChild(routes)],
    exports: [RouterModule],
})
export class UserListRoutingModule { }
