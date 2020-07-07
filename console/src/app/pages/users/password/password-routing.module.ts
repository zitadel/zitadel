import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';

import { PasswordComponent } from '../password/password.component';

const routes: Routes = [
    {
        path: '',
        component: PasswordComponent,
    },
];

@NgModule({
    imports: [RouterModule.forChild(routes)],
    exports: [RouterModule],
})
export class PasswordRoutingModule { }
