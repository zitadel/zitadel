import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';

import { IamMembersComponent } from './iam-members.component';

const routes: Routes = [
    {
        path: '',
        component: IamMembersComponent,
        data: { animation: 'AddPage' },
    },
];

@NgModule({
    imports: [RouterModule.forChild(routes)],
    exports: [RouterModule],
})
export class IamMembersRoutingModule { }
