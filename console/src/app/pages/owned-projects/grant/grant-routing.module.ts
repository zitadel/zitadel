import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';

import { GrantComponent } from './grant.component';

const routes: Routes = [
    {
        path: '',
        component: GrantComponent,
        data: { animation: 'AddPage' },
    },
];

@NgModule({
    imports: [RouterModule.forChild(routes)],
    exports: [RouterModule],
})
export class GrantRoutingModule { }
