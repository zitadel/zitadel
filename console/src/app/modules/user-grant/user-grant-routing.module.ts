import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';

import { UserGrantComponent } from './user-grant.component';

const routes: Routes = [
    {
        path: '',
        component: UserGrantComponent,
        data: { animation: 'AddPage' },
    },
];

@NgModule({
    imports: [RouterModule.forChild(routes)],
    exports: [RouterModule],
})
export class UserGrantRoutingModule { }
