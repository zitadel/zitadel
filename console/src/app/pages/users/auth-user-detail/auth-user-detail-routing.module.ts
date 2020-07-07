import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';

import { AuthUserDetailComponent } from './auth-user-detail.component';

const routes: Routes = [
    {
        path: '',
        component: AuthUserDetailComponent,
    },
];

@NgModule({
    imports: [RouterModule.forChild(routes)],
    exports: [RouterModule],
})
export class AuthUserDetailRoutingModule { }
