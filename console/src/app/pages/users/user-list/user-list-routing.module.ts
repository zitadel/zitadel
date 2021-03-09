import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';
import { Type } from 'src/app/proto/generated/zitadel/user_pb';

import { UserListComponent } from './user-list.component';


const routes: Routes = [
    {
        path: 'humans',
        component: UserListComponent,
        data: {
            animation: 'HomePage',
            type: Type.TYPE_HUMAN,
        },
    },
    {
        path: 'machines',
        component: UserListComponent,
        data: {
            animation: 'HomePage',
            type: Type.TYPE_MACHINE,
        },
    },
];

@NgModule({
    imports: [RouterModule.forChild(routes)],
    exports: [RouterModule],
})
export class UserListRoutingModule { }
