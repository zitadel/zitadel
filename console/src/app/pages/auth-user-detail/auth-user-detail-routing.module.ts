import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';
import { AuthGuard } from 'src/app/guards/auth.guard';

import { AuthUserDetailComponent } from './auth-user-detail.component';

const routes: Routes = [
    {
        path: 'me',
        component: AuthUserDetailComponent,
    },
    {
        path: 'me/password',
        loadChildren: () => import('src/app/pages/password/password.module').then(m => m.PasswordModule),
        canActivate: [AuthGuard],
    },
];

@NgModule({
    imports: [RouterModule.forChild(routes)],
    exports: [RouterModule],
})
export class AuthUserDetailRoutingModule { }
