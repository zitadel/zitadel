import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';

import { GrantsComponent } from './grants.component';

const routes: Routes = [
    {
        path: '',
        component: GrantsComponent,
    },
];

@NgModule({
    imports: [RouterModule.forChild(routes)],
    exports: [RouterModule],
})
export class GrantsRoutingModule { }
