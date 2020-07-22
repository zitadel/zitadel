import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';

import { UserListComponent } from './user-list.component';


const routes: Routes = [
    {
        path: '',
        component: UserListComponent,
        data: { animation: 'HomePage' },
    },
];

@NgModule({
    imports: [RouterModule.forChild(routes)],
    exports: [RouterModule],
})
export class UserListRoutingModule { }
