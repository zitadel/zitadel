import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';

import { IamComponent } from './iam.component';

const routes: Routes = [
    {
        path: '',
        component: IamComponent,
    },
];

@NgModule({
    imports: [RouterModule.forChild(routes)],
    exports: [RouterModule],
})
export class IamRoutingModule { }
