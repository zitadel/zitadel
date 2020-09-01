import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';

import { UserListComponent, UserType } from './user-list.component';


const routes: Routes = [
    {
        path: 'humans',
        component: UserListComponent,
        data: {
            animation: 'HomePage',
            type: UserType.HUMAN,
        },
    },
    {
        path: 'machines',
        component: UserListComponent,
        data: {
            animation: 'HomePage',
            type: UserType.MACHINE,
        },
    },
];

@NgModule({
    imports: [RouterModule.forChild(routes)],
    exports: [RouterModule],
})
export class UserListRoutingModule { }
